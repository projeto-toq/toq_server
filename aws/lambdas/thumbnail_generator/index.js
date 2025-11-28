const { S3Client, GetObjectCommand, PutObjectCommand } = require('@aws-sdk/client-s3');
const sharp = require('sharp');
const { Readable } = require('stream');

const s3Client = new S3Client({ region: process.env.AWS_REGION || 'us-east-1' });
const MEDIA_BUCKET = process.env.MEDIA_BUCKET || 'toq-listing-medias';

// Thumbnail sizes
const THUMBNAIL_SIZES = [
    { name: 'small', width: 320, height: 240 },
    { name: 'medium', width: 640, height: 480 },
    { name: 'large', width: 1280, height: 960 }
];

exports.handler = async (event) => {
    console.log('Thumbnails Lambda - Event received:', JSON.stringify(event, null, 2));
    
    const { batchId, listingId, validAssets, traceparent } = event.body || event;
    
    if (!Array.isArray(validAssets)) {
        throw new Error('validAssets array is required');
    }
    
    // Filter only photo assets
    const photoAssets = validAssets.filter(asset => 
        asset.assetType && asset.assetType.startsWith('PHOTO_')
    );
    
    if (photoAssets.length === 0) {
        console.log('No photo assets to process');
        return {
            statusCode: 200,
            body: {
                status: 'no_photos_to_process',
                batchId,
                listingId,
                traceparent,
                thumbnails: []
            }
        };
    }
    
    const thumbnailResults = [];
    const errors = [];
    
    for (const asset of photoAssets) {
        try {
            const thumbs = await generateThumbnails(asset, listingId);
            thumbnailResults.push(...thumbs);
        } catch (error) {
            console.error(`Error generating thumbnails for ${asset.sourceKey}:`, error);
            errors.push({
                sourceKey: asset.sourceKey,
                error: error.message
            });
        }
    }
    
    return {
        statusCode: 200,
        body: {
            status: 'thumbnails_generated',
            batchId,
            listingId,
            traceparent,
            assetsProcessed: photoAssets.length,
            thumbnailsGenerated: thumbnailResults.length,
            thumbnails: thumbnailResults,
            errors: errors.length > 0 ? errors : undefined
        }
    };
};

async function generateThumbnails(asset, listingId) {
    const { sourceKey } = asset;
    
    // Download original image
    const getCommand = new GetObjectCommand({
        Bucket: MEDIA_BUCKET,
        Key: sourceKey
    });
    
    const response = await s3Client.send(getCommand);
    const imageBuffer = await streamToBuffer(response.Body);
    
    // Get image metadata
    const metadata = await sharp(imageBuffer).metadata();
    console.log(`Processing image ${sourceKey}: ${metadata.width}x${metadata.height}, format: ${metadata.format}`);
    
    const thumbnails = [];
    
    // Generate thumbnails for each size
    for (const size of THUMBNAIL_SIZES) {
        try {
            const thumbnailBuffer = await sharp(imageBuffer)
                .resize(size.width, size.height, {
                    fit: 'inside',
                    withoutEnlargement: true
                })
                .jpeg({
                    quality: 85,
                    progressive: true
                })
                .toBuffer();
            
            // Generate thumbnail key
            const pathParts = sourceKey.split('/');
            const filename = pathParts[pathParts.length - 1];
            const nameWithoutExt = filename.replace(/\.[^/.]+$/, '');
            const thumbnailKey = `processed/thumb/${listingId}/${size.name}/${nameWithoutExt}_${size.width}x${size.height}.jpg`;
            
            // Upload thumbnail
            await s3Client.send(new PutObjectCommand({
                Bucket: MEDIA_BUCKET,
                Key: thumbnailKey,
                Body: thumbnailBuffer,
                ContentType: 'image/jpeg',
                Metadata: {
                    'original-key': sourceKey,
                    'thumbnail-size': size.name,
                    'listing-id': listingId.toString()
                }
            }));
            
            thumbnails.push({
                originalKey: sourceKey,
                thumbnailKey,
                size: size.name,
                width: size.width,
                height: size.height,
                bytes: thumbnailBuffer.length
            });
            
            console.log(`Generated thumbnail: ${thumbnailKey} (${thumbnailBuffer.length} bytes)`);
            
        } catch (error) {
            console.error(`Failed to generate ${size.name} thumbnail for ${sourceKey}:`, error);
        }
    }
    
    return thumbnails;
}

async function streamToBuffer(stream) {
    if (stream instanceof Buffer) {
        return stream;
    }
    
    const chunks = [];
    for await (const chunk of stream) {
        chunks.push(chunk);
    }
    return Buffer.concat(chunks);
}
