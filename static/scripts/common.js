// common.js

// Handle file download click
async function fileClick(filename) {
    try {
        // Trigger file download by opening the download URL
        const response = await fetch(`/download/${filename}`);
        if (!response.ok) {
            throw new Error(`Failed to download file: ${filename}`);
        }
        // Open the file content in the browser
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.style.display = 'none';
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
    } catch (error) {
        console.error("Download failed:", error);
    }
}

// Handle file upload to Slack
function slackFileClick(filename, webhookUrl) {
    // Implement Slack file upload logic here (could involve API calls)
    console.log("Uploading file to Slack:", filename);
}
