async function downloadFileToBlob(fileUrl) {
    try {
        // 1. Download the file
        const response = await fetch(fileUrl);

        if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
        }

        // Get the blob of the file
        const blob = await response.blob();
        return blob
    }
    catch (error) {
        console.error('Error downloading or posting file:', error);
        throw error; // Re-throw the error to be handled by the caller if needed
    }
}