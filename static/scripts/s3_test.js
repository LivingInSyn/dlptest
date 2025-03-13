async function s3fileClick(filename) {
    let dlurl = "/static/downloads/"+filename;
    let presignUrl = "/generateS3Token?filename="+filename

    try {
        const response = await fetch(presignUrl);
        if (!response.ok) {
          throw new Error(`HTTP error getting presigned URL! status: ${response.status}`);
        }
        const presignData = await response.json();
        let result = await downloadAndPostFileS3(dlurl, presignData);
    } catch (error) {
        console.error("Fetching JSON failed:", error);
    }


    const presignReq = await fetch(presignUrl)
    if (!Response.ok) {
        throw new Error(`HTTP error getting presign URL! status: ${response.status}`);
    }
}

async function downloadAndPostFileS3(fileUrl, presignData){
    let filename = fileUrl.substring(fileUrl.lastIndexOf('/') + 1);
    if (!filename) { // If no filename found in URL, generate a generic one
        filename = "downloaded_file";
    }
    id_to_update = filename + "-status-s3"
    status_span = document.getElementById(id_to_update);
    try {
         // 1. Download the file
         const response = await fetch(fileUrl);

         if (!response.ok) {
         throw new Error(`HTTP error fetching file in s3 test! status: ${response.status}`);
         }
 
         // Get the blob of the file
         const blob = await downloadFileToBlob(fileUrl)

         //
         const uploadResponse = await fetch(presignData.url, {
            method: "PUT",
            body: blob,
            //headers: { "Content-Type": file.type }
        })

        if (!uploadResponse.ok) {
            const errorText = await uploadResponse.text(); // Try to get error message from server
            //throw new Error(`POST error! status: ${postResponse.status}, message: ${errorText}`);
            status_span.textContent = "PUT to s3 failed - status: " + postResponse.status;
            return
        }

        console.log('File PUT to s3 successfully:');
        status_span.textContent = "PUT successful";

        return 0; // Return the server's response


    } catch (error) {
        status_span.textContent = "POST not successful" + error;
        console.error('Error downloading or posting file:', error);
        throw error; // Re-throw the error to be handled by the caller if needed
    }
}