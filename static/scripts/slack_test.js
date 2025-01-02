async function slackFileClick(filename, webhook) {
    let dlurl = "http://localhost:8080/static/downloads/"+filename;
    let posturl = "http://localhost:8080/upload";
    const blob = downloadFileToBlob(dlurl)
    postFileToSlackWebhook(blob, filename, webhook, "DLP Test")
}

async function postFileToSlackWebhook(file, filename, webhookUrl, initialComment = "") {
    try {
      if (!file || !(file instanceof Blob) ) {
          throw new Error("File must be a Blob object.");
      }
  
      if (!filename || typeof filename !== 'string'){
          throw new Error("Filename must be a non-empty string.");
      }
  
      if (!webhookUrl || typeof webhookUrl !== 'string') {
          throw new Error("Webhook URL must be a non-empty string.");
      }
  
      const formData = new FormData();
      formData.append('file', file, filename);
      if (initialComment) {
        formData.append('initial_comment', initialComment);
      }
  
  
      const response = await fetch(webhookUrl, {
        method: 'POST',
        body: formData,
      });
  
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`Slack webhook POST error: ${response.status} - ${errorText}`);
      }
  
      const result = await response.json(); // Or response.text() if Slack doesn't return JSON
      console.log('File posted to Slack:', result);
      return result;
  
    } catch (error) {
      console.error('Error posting file to Slack:', error);
      throw error; // Re-throw for handling by the caller
    }
  }

  async function testSlackPost() {
    const testWebhookUrl = "YOUR_SLACK_WEBHOOK_URL"; // Replace with your actual webhook URL
    if (testWebhookUrl == "YOUR_SLACK_WEBHOOK_URL"){
        console.log("Please replace YOUR_SLACK_WEBHOOK_URL with a valid webhook url to test this function");
        return;
    }

    try {
        // Create a test Blob (you would normally get this from a file input or other source)
        const testText = "This is a test file for Slack.";
        const testBlob = new Blob([testText], { type: 'text/plain' });
        const result = await postFileToSlackWebhook(testBlob, "test_file.txt", testWebhookUrl, "Testing file upload");
        console.log("Result from slack", result)
    } catch (error) {
        console.error("Test failed:", error);
    }

    try{
        await postFileToSlackWebhook("not a blob", "test_file.txt", testWebhookUrl, "Testing file upload");
    } catch (error) {
        console.log("Correctly caught error for non blob", error);
    }

    try{
        const testText = "This is a test file for Slack.";
        const testBlob = new Blob([testText], { type: 'text/plain' });
        await postFileToSlackWebhook(testBlob, 123, testWebhookUrl, "Testing file upload");
    } catch (error) {
        console.log("Correctly caught error for non string filename", error);
    }

    try{
        const testText = "This is a test file for Slack.";
        const testBlob = new Blob([testText], { type: 'text/plain' });
        await postFileToSlackWebhook(testBlob, "test", 123, "Testing file upload");
    } catch (error) {
        console.log("Correctly caught error for non string webhook url", error);
    }
}