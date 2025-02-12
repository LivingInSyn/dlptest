// Function to download a file and return it as a Blob
async function downloadFileToBlob(fileUrl) {
    try {
        console.log(`Downloading file from: ${fileUrl}`);

        const response = await fetch(fileUrl);

        if (!response.ok) {
            throw new Error(`Failed to download file. HTTP Status: ${response.status} - ${response.statusText}`);
        }

        // Convert to Blob
        const blob = await response.blob();
        console.log("File successfully downloaded as Blob.");
        return blob;
    } catch (error) {
        console.error("Error downloading file:", error.message);
        throw error;
    }
}

// Function to create FormData with pre-signed data for S3 uploads
function createFormData(file, preSignedData) {
    if (!file || !preSignedData) {
        throw new Error('Both file and pre-signed data must be provided.');
    }

    const formData = new FormData();
    formData.append('file', file);
    formData.append('key', preSignedData.fields.key);
    formData.append('AWSAccessKeyId', preSignedData.fields.AWSAccessKeyId); // Ensure this is dynamically set
    formData.append('policy', preSignedData.fields.policy);
    formData.append('signature', preSignedData.fields.signature);

    return formData;
}
