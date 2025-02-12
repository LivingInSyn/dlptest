// http_test.js

// Function to handle the file click and process the download/upload
async function fileClick(filename) {
    const dlurl = `http://localhost:8080/static/downloads/${filename}`;
    const posturl = "http://localhost:8080/upload";
    
    try {
        let result = await downloadAndPostFile(dlurl, posturl);
        console.log("Result from server:", result);
    } catch (error) {
        console.error("Error during file download and posting:", error);
    }
}

// Function to download and post the file
async function downloadAndPostFile(fileUrl, postUrl) {
    const filename = fileUrl.substring(fileUrl.lastIndexOf('/') + 1);
    const status_span = document.getElementById(`${filename}-status`);

    try {
        // Download the file as a Blob
        const blob = await downloadFileToBlob(fileUrl);

        // Create FormData for the POST request
        const formData = new FormData();
        formData.append('file', blob, filename);

        // Send the file to the server via POST request
        const postResponse = await fetch(postUrl, {
            method: 'POST',
            body: formData,
        });

        if (!postResponse.ok) {
            const errorText = await postResponse.text();
            if (status_span) {
                status_span.textContent = `POST failed - Status: ${postResponse.status}`;
            }
            throw new Error(`POST error: ${postResponse.status}, Message: ${errorText}`);
        }

        if (status_span) {
            status_span.textContent = "POST successful";
        }

        console.log('File posted successfully');
        return await postResponse.json();
    } catch (error) {
        if (status_span) {
            status_span.textContent = `POST not successful - ${error.message}`;
        }
        console.error('Error downloading or posting file:', error);
        throw error;
    }
}
