/* slack_test.js */

async function slackFileClick(filename, webhook) {
    const dlurl = `http://localhost:8080/static/downloads/${filename}`;
    const statusElement = document.getElementById(`${filename}-status`);

    try {
        statusElement.textContent = "Downloading file...";
        const blob = await downloadFileToBlob(dlurl);

        statusElement.textContent = "Uploading file to Slack...";
        await uploadFileToSlack(blob, filename, webhook);

        statusElement.textContent = "File uploaded successfully!";
    } catch (error) {
        statusElement.textContent = "Error during upload.";
        console.error("Error during file click process:", error);
    }
}

async function uploadFileToSlack(file, filename, webhook) {
    const formData = new FormData();
    formData.append("file", file, filename);

    try {
        const response = await fetch(webhook, {
            method: "POST",
            body: formData
        });

        const result = await response.json();
        if (!result.ok) throw new Error(result.error);

        console.log("File posted to Slack:", result);
    } catch (error) {
        console.error("Error posting file to Slack:", error);
    }
}