import json
import os
import boto3
import zipfile
from io import BytesIO
from datetime import datetime

s3_client = boto3.client('s3', region_name=os.environ.get('AWS_REGION', 'us-east-1'))
MEDIA_BUCKET = os.environ.get('MEDIA_BUCKET', 'toq-listing-medias')

def lambda_handler(event, context):
    print(f'ZIP Lambda - Event received: {json.dumps(event)}')
    
    body = event.get('body', event)
    batch_id = body.get('batchId')
    listing_id = body.get('listingId')
    thumbnails = body.get('thumbnails', [])
    valid_assets = body.get('validAssets', [])
    traceparent = body.get('traceparent')
    
    if not batch_id or not listing_id:
        raise ValueError('batchId and listingId are required')
    
    # Collect all assets to zip
    assets_to_zip = []
    
    # Add original assets
    for asset in valid_assets:
        assets_to_zip.append({
            'key': asset['sourceKey'],
            'type': 'original',
            'asset_type': asset.get('assetType', 'UNKNOWN')
        })
    
    # Add thumbnails
    for thumb in thumbnails:
        assets_to_zip.append({
            'key': thumb['thumbnailKey'],
            'type': 'thumbnail',
            'size': thumb['size']
        })
    
    if len(assets_to_zip) == 0:
        print('No assets to zip')
        return {
            'statusCode': 200,
            'body': {
                'status': 'no_assets_to_zip',
                'batchId': batch_id,
                'listingId': listing_id,
                'traceparent': traceparent
            }
        }
    
    # Create ZIP bundles
    try:
        zip_results = create_zip_bundles(listing_id, batch_id, assets_to_zip)
        
        return {
            'statusCode': 200,
            'body': {
                'status': 'zip_created',
                'batchId': batch_id,
                'listingId': listing_id,
                'traceparent': traceparent,
                'assetsZipped': len(assets_to_zip),
                'zipBundles': zip_results
            }
        }
        
    except Exception as error:
        print(f'Error creating ZIP: {error}')
        raise

def create_zip_bundles(listing_id, batch_id, assets):
    """Create ZIP bundles - one complete, one originals-only, one thumbs-only"""
    
    results = []
    timestamp = datetime.utcnow().strftime('%Y%m%d_%H%M%S')
    
    # Bundle 1: Complete (all assets)
    complete_zip_key = f"zip/{listing_id}/complete_{batch_id}_{timestamp}.zip"
    complete_size = create_zip(complete_zip_key, assets, filter_type=None)
    results.append({
        'bundleType': 'complete',
        'zipKey': complete_zip_key,
        'sizeBytes': complete_size,
        'assetsCount': len(assets)
    })
    
    # Bundle 2: Originals only
    originals = [a for a in assets if a['type'] == 'original']
    if len(originals) > 0:
        originals_zip_key = f"zip/{listing_id}/originals_{batch_id}_{timestamp}.zip"
        originals_size = create_zip(originals_zip_key, originals, filter_type='original')
        results.append({
            'bundleType': 'originals',
            'zipKey': originals_zip_key,
            'sizeBytes': originals_size,
            'assetsCount': len(originals)
        })
    
    # Bundle 3: Thumbnails only
    thumbnails = [a for a in assets if a['type'] == 'thumbnail']
    if len(thumbnails) > 0:
        thumbs_zip_key = f"zip/{listing_id}/thumbnails_{batch_id}_{timestamp}.zip"
        thumbs_size = create_zip(thumbs_zip_key, thumbnails, filter_type='thumbnail')
        results.append({
            'bundleType': 'thumbnails',
            'zipKey': thumbs_zip_key,
            'sizeBytes': thumbs_size,
            'assetsCount': len(thumbnails)
        })
    
    return results

def create_zip(zip_key, assets, filter_type=None):
    """Create a ZIP file in memory and upload to S3"""
    
    zip_buffer = BytesIO()
    total_size = 0
    
    with zipfile.ZipFile(zip_buffer, 'w', zipfile.ZIP_DEFLATED) as zip_file:
        for asset in assets:
            if filter_type and asset['type'] != filter_type:
                continue
            
            try:
                # Download asset from S3
                response = s3_client.get_object(Bucket=MEDIA_BUCKET, Key=asset['key'])
                file_data = response['Body'].read()
                
                # Add to ZIP with organized structure
                arcname = get_archive_name(asset)
                zip_file.writestr(arcname, file_data)
                total_size += len(file_data)
                
                print(f"Added to ZIP: {arcname} ({len(file_data)} bytes)")
                
            except Exception as error:
                print(f"Error adding {asset['key']} to ZIP: {error}")
                # Continue with other files
    
    # Upload ZIP to S3
    zip_buffer.seek(0)
    s3_client.put_object(
        Bucket=MEDIA_BUCKET,
        Key=zip_key,
        Body=zip_buffer.getvalue(),
        ContentType='application/zip',
        Metadata={
            'batch-type': filter_type or 'complete',
            'assets-count': str(len(assets))
        }
    )
    
    zip_size = len(zip_buffer.getvalue())
    print(f"ZIP created: {zip_key} ({zip_size} bytes, {len(assets)} files)")
    
    return zip_size

def get_archive_name(asset):
    """Generate organized archive name based on asset type"""
    
    key = asset['key']
    filename = key.split('/')[-1]
    
    if asset['type'] == 'thumbnail':
        size = asset.get('size', 'unknown')
        return f"thumbnails/{size}/{filename}"
    elif asset['type'] == 'original':
        asset_type = asset.get('asset_type', 'UNKNOWN')
        if 'PHOTO' in asset_type:
            return f"photos/{filename}"
        elif 'VIDEO' in asset_type:
            return f"videos/{filename}"
        elif 'PROJECT' in asset_type:
            return f"projects/{filename}"
        else:
            return f"originals/{filename}"
    
    return filename
