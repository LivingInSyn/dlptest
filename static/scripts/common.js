// common.js

async function downloadFileToBlob(fileUrl) {
    try {
        // 1. Download the file
        const response = await fetch(fileUrl);

        if (!response.ok) {
            throw new Error(`Failed to fetch the file from ${fileUrl}. HTTP error! Status: ${response.status}`);
        }

        // 2. Get the blob of the file
        const blob = await response.blob();
        return blob;
    }
    catch (error) {
        console.error(`Error downloading file from ${fileUrl}:`, error);
        throw error; // Re-throw the error to be handled by the caller if needed
    }
}
