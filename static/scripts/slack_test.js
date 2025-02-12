/*slack_test.js*/

async function slackFileClick(filename, webhook) {
    const dlurl = `http://localhost:8080/static/downloads/${filename}`;
    const statusElement = document.getElementById("status");

    try {
        statusElement.textContent = "Downloading file...";

        const blob = await downloadFileToBlob(dlurl);

        statusElement.textContent = "Uploading file to Slack...";

        await uploadFileToSlack(blob, filename);

        statusElement.textContent = "File uploaded successfully!";
    } catch (error) {
        statusElement.textContent = "Error during upload.";
        console.error("Error during file click process:", error);
    }
}

async function uploadFileToSlack(file, filename) {
    const slackToken = "YOUR_SLACK_BOT_TOKEN"; // Replace with your Slack token
    const slackChannel = "YOUR_SLACK_CHANNEL_ID"; // Replace with your Slack channel ID

    const formData = new FormData();
    formData.append("file", file, filename);
    formData.append("channels", slackChannel);

    try {
        const response = await fetch("https://slack.com/api/files.upload", {
            method: "POST",
            headers: { Authorization: `Bearer ${slackToken}` },
            body: formData
        });

        const result = await response.json();
        if (!result.ok) throw new Error(result.error);

        console.log("File posted to Slack:", result);
    } catch (error) {
        console.error("Error posting file to Slack:", error);
    }
}
