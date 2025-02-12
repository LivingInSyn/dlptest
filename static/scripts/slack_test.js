// Function to initiate the file click process for Slack upload
async function slackFileClick(filename, webhook) {
    const dlurl = `http://localhost:8080/static/downloads/${filename}`;
    const statusElement = document.getElementById("status");

    try {
        statusElement.textContent = "Downloading file...";

        // Download the file as a blob
        const blob = await downloadFileToBlob(dlurl);

        statusElement.textContent = "Uploading file to Slack...";

        // Upload file to Slack via webhook
        await postFileToSlackWebhook(blob, filename, webhook);

        statusElement.textContent = "File uploaded successfully!";
    } catch (error) {
        statusElement.textContent = "Error during upload.";
        console.error("Error during file click process:", error);
    }
}

// Function to post a file to Slack via webhook
async function postFileToSlackWebhook(file, filename, webhookUrl, initialComment = "") {
    if (!(file instanceof Blob)) {
        throw new Error("File must be a Blob object.");
    }

    if (typeof filename !== 'string' || !filename.trim()) {
        throw new Error("Filename must be a non-empty string.");
    }

    if (typeof webhookUrl !== 'string' || !webhookUrl.trim()) {
        throw new Error("Webhook URL must be a non-empty string.");
    }

    const formData = new FormData();
    formData.append('file', file, filename);
    if (initialComment) {
        formData.append('initial_comment', initialComment);
    }

    try {
        const response = await fetch(webhookUrl, {
            method: 'POST',
            body: formData,
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Slack webhook POST error: ${response.status} - ${errorText}`);
        }

        const result = await response.json();
        console.log('File posted to Slack:', result);
        return result;
    } catch (error) {
        console.error('Error posting file to Slack:', error);
        throw error;
    }
}

// Helper function to download file as a Blob
async function downloadFileToBlob(fileUrl) {
    const response = await fetch(fileUrl);

    if (!response.ok) {
        throw new Error(`Failed to fetch file from ${fileUrl}. HTTP error: ${response.status}`);
    }

    return response.blob();
}
