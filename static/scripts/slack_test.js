// slack_tests.js

// Function to initiate the file click process
async function slackFileClick(filename, webhook) {
    let dlurl = `http://localhost:8080/static/downloads/${filename}`;
    let posturl = "http://localhost:8080/upload";
    
    try {
        // Download the file as a blob
        const blob = await downloadFileToBlob(dlurl);
        // Post the file to Slack
        await postFileToSlackWebhook(blob, filename, webhook, "DLP Test");
    } catch (error) {
        console.error("Error during file click process:", error);
    }
}

// Function to post file to Slack via webhook
async function postFileToSlackWebhook(file, filename, webhookUrl, initialComment = "") {
    try {
        // Validate parameters
        if (!(file instanceof Blob)) {
            throw new Error("File must be a Blob object.");
        }

        if (typeof filename !== 'string' || !filename.trim()) {
            throw new Error("Filename must be a non-empty string.");
        }

        if (typeof webhookUrl !== 'string' || !webhookUrl.trim()) {
            throw new Error("Webhook URL must be a non-empty string.");
        }

        // Create FormData object and append file and comment
        const formData = new FormData();
        formData.append('file', file, filename);
        if (initialComment) {
            formData.append('initial_comment', initialComment);
        }

        // Send the POST request to Slack webhook
        const response = await fetch(webhookUrl, {
            method: 'POST',
            body: formData,
        });

        // Handle unsuccessful response
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Slack webhook POST error: ${response.status} - ${errorText}`);
        }

        // Parse the response from Slack (assuming JSON)
        const result = await response.json();
        console.log('File posted to Slack:', result);
        return result;

    } catch (error) {
        console.error('Error posting file to Slack:', error);
        throw error; // Re-throw for handling by the caller
    }
}

// Test function for Slack file posting
async function testSlackPost() {
    const testWebhookUrl = "YOUR_SLACK_WEBHOOK_URL"; // Replace with your actual webhook URL
    
    if (testWebhookUrl === "YOUR_SLACK_WEBHOOK_URL") {
        console.log("Please replace YOUR_SLACK_WEBHOOK_URL with a valid webhook URL to test this function");
        return;
    }

    try {
        // Create a test Blob (this would normally be fetched from a file input or other source)
        const testText = "This is a test file for Slack.";
        const testBlob = new Blob([testText], { type: 'text/plain' });
        const result = await postFileToSlackWebhook(testBlob, "test_file.txt", testWebhookUrl, "Testing file upload");
        console.log("Result from Slack:", result);
    } catch (error) {
        console.error("Test failed:", error);
    }

    try {
        await postFileToSlackWebhook("not a blob", "test_file.txt", testWebhookUrl, "Testing file upload");
    } catch (error) {
        console.log("Correctly caught error for non-blob:", error);
    }

    try {
        const testText = "This is a test file for Slack.";
        const testBlob = new Blob([testText], { type: 'text/plain' });
        await postFileToSlackWebhook(testBlob, 123, testWebhookUrl, "Testing file upload");
    } catch (error) {
        console.log("Correctly caught error for non-string filename:", error);
    }

    try {
        const testText = "This is a test file for Slack.";
        const testBlob = new Blob([testText], { type: 'text/plain' });
        await postFileToSlackWebhook(testBlob, "test", 123, "Testing file upload");
    } catch (error) {
        console.log("Correctly caught error for non-string webhook URL:", error);
    }
}
