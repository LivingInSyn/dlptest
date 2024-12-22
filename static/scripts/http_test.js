async function fileClick(filename) {
    let dlurl = "http://localhost:8080/static/downloads/"+filename;
    let posturl = "http://localhost:8080/upload";
    let result = await downloadAndPostFile(dlurl, posturl);
    console.log("Result from server:", result)

}
async function downloadAndPostFile(fileUrl, postUrl) {
    let filename = fileUrl.substring(fileUrl.lastIndexOf('/') + 1);
    if (!filename) { // If no filename found in URL, generate a generic one
        filename = "downloaded_file";
    }
    id_to_update = filename + "-status"
    status_span = document.getElementById(id_to_update);
    try {
        // 1. Download the file
        const response = await fetch(fileUrl);

        if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
        }

        // Get the blob of the file
        const blob = await response.blob();

        // 2. Create a FormData object for the POST request
        const formData = new FormData();
        
        // Append the blob to the FormData. You can specify a custom filename.
        formData.append('file', blob, filename); // 'file' is the field name expected by the server

        // 3. POST the file
        const postResponse = await fetch(postUrl, {
        method: 'POST',
        body: formData,
        });

        if (!postResponse.ok) {
            const errorText = await postResponse.text(); // Try to get error message from server
            //throw new Error(`POST error! status: ${postResponse.status}, message: ${errorText}`);
            status_span.textContent = "POST failed - status: " + postResponse.status;
            return
        }

        console.log('File posted successfully:');
        status_span.textContent = "POST successful";

        return 0; // Return the server's response
    } catch (error) {
        status_span.textContent = "POST not successful" + error;
        console.error('Error downloading or posting file:', error);
        throw error; // Re-throw the error to be handled by the caller if needed
    }
}