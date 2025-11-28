const { S3Client, HeadObjectCommand, PutObjectCommand } = require('@aws-sdk/client-s3');
const { SFNClient, StartExecutionCommand } = require("@aws-sdk/client-sfn");

const s3Client = new S3Client({ region: process.env.AWS_REGION || 'us-east-1' });
const sfnClient = new SFNClient({ region: process.env.AWS_REGION || 'us-east-1' });

const MEDIA_BUCKET = process.env.MEDIA_BUCKET || 'toq-listing-medias';
const STATE_MACHINE_ARN = process.env.STATE_MACHINE_ARN || 'arn:aws:states:us-east-1:058264253741:stateMachine:listing-media-processing-sm-staging';

async function logToS3(key, data) {
    try {
        await s3Client.send(new PutObjectCommand({
            Bucket: MEDIA_BUCKET,
            Key: `debug/${key}`,
            Body: JSON.stringify(data, null, 2)
        }));
    } catch (e) {
        console.error("Failed to log to S3", e);
    }
}

exports.handler = async (event) => {
    console.log('Validate Lambda - Event received:', JSON.stringify(event, null, 2));
    await logToS3(`event-${Date.now()}.json`, event);

    // Handle SQS Trigger
    if (event.Records && Array.isArray(event.Records)) {
        const results = [];
        for (const record of event.Records) {
            try {
                const body = JSON.parse(record.body);
                await logToS3(`sqs-body-${Date.now()}.json`, body);

                // 1. Validate (Fail fast)
                const validationResult = await processEvent(body);

                // 2. Start Step Function if valid
                if (validationResult.status === 'validated') {
                    await startStepFunction(body, validationResult.traceparent);
                } else {
                    console.warn("Validation failed, skipping Step Function execution", JSON.stringify(validationResult.errors));
                }

                results.push(validationResult);
            } catch (err) {
                console.error("Failed to process record", record.messageId, err);
                await logToS3(`error-${Date.now()}.json`, { error: err.message, stack: err.stack });
                throw err;
            }
        }
        return results;
    }

    // Handle Direct Invocation (e.g. from Step Function)
    // Just validate and return result. DO NOT start Step Function.
    return await processEvent(event);
}; async function startStepFunction(payload, traceparent) {
    const { batchId } = payload;
    try {
        const sfnCommand = new StartExecutionCommand({
            stateMachineArn: STATE_MACHINE_ARN,
            name: `batch-${batchId}-${Date.now()}`,
            input: JSON.stringify(payload), // Pass original payload so SF can use it
            traceHeader: traceparent
        });

        const sfnResult = await sfnClient.send(sfnCommand);
        console.log(`Started Step Function execution: ${sfnResult.executionArn}`);
        return sfnResult;
    } catch (sfnError) {
        console.error("Failed to start Step Function:", sfnError);
        throw sfnError;
    }
}

async function processEvent(payload) {
    const { batchId, listingId, assets, traceparent } = payload;

    if (!batchId || !listingId || !Array.isArray(assets)) {
        throw new Error('Invalid input: batchId, listingId, and assets array are required');
    }

    const validationResults = [];
    const errors = [];

    for (const assetKey of assets) {
        try {
            const result = await validateAsset(assetKey);
            validationResults.push(result);
        } catch (error) {
            console.error(`Error validating asset ${assetKey}:`, error);
            errors.push({
                sourceKey: assetKey,
                error: error.message
            });
        }
    }

    const hasErrors = errors.length > 0;
    const outputPayload = {
        status: hasErrors ? 'validation_failed' : 'validated',
        batchId,
        listingId,
        traceparent,
        assetsValidated: validationResults.length,
        validAssets: validationResults,
        errors: hasErrors ? errors : undefined,
        hasVideos: validationResults.some(a => a.assetType && a.assetType.includes('VIDEO'))
    };

    return outputPayload;
}

async function validateAsset(assetKey) {
    // Basic validation: check if object exists in S3
    // The 'asset' input from SQS is just the key string based on the logs: 
    // "assets": ["51/raw/photo/vertical/2025-11-27/photo-001.jpg", ...]

    try {
        const command = new HeadObjectCommand({
            Bucket: MEDIA_BUCKET,
            Key: assetKey
        });

        const response = await s3Client.send(command);

        // Determine type based on key path or extension (simplified logic)
        let assetType = 'PHOTO';
        if (assetKey.includes('/video/')) assetType = 'VIDEO';

        return {
            rawKey: assetKey,
            sourceKey: assetKey,
            assetType: assetType,
            size: response.ContentLength,
            contentType: response.ContentType,
            etag: response.ETag
        };
    } catch (error) {
        throw new Error(`S3 HeadObject failed for ${assetKey}: ${error.message}`);
    }
}
