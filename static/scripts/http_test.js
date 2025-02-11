// http_test.js

async function fileClick(filename) {
    let dlurl = `http://localhost:8080/static/downloads/${filename}`;
    let posturl = "http://localhost:8080/upload";
    
    try {
        let result = await downloadAndPostFile(dlurl, posturl);
        console.log("Result from server:", result);
    } catch (error) {
        console.error("Error during file download and posting:", error);
    }
}

async function downloadAndPostFile(fileUrl, postUrl) {
    // Extract filename from URL
    let filename = fileUrl.substring(fileUrl.lastIndexOf('/') + 1);
    if (!filename) { // If no filename found in URL, generate a default name
        filename = "downloaded_file";
    }

    // Check if the status span exists before updating
    const id_to_update = `${filename}-status`;
    const status_span = document.getElementById(id_to_update);
    
    try {
        // 1. Download the file
        const response = await fetch(fileUrl);

        if (!response.ok) {
            throw new Error(`Failed to fetch file from ${fileUrl}. HTTP error: ${response.status}`);
        }

        // Get the blob of the file
        const blob = await downloadFileToBlob(fileUrl);

        // 2. Create a FormData object for the POST request
        const formData = new FormData();
        formData.append('file', blob, filename); // 'file' is the expected field name by the server

        // 3. POST the file
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
        return await postResponse.json(); // Assuming the server returns a JSON response
    } catch (error) {
        if (status_span) {
            status_span.textContent = `POST not successful - ${error.message}`;
        }
        console.error('Error downloading or posting file:', error);
        throw error; // Re-throw for handling by the caller
    }
}
