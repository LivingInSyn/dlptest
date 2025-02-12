async function downloadFileToBlob(fileUrl) {
    try {
        console.log(`Downloading file from: ${fileUrl}`);

        // Fetch the file
        const response = await fetch(fileUrl);

        if (!response.ok) {
            throw new Error(`Failed to download file. HTTP Status: ${response.status} - ${response.statusText}`);
        }

        // Convert to blob
        const blob = await response.blob();
        console.log("File successfully downloaded as Blob.");
        return blob;
    } catch (error) {
        console.error("Error downloading file:", error.message);
        throw error; // Ensure error propagates for handling upstream
    }
}
