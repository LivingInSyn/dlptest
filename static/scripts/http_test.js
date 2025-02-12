/*http_test.js*/


async function fileClick(filename, uploadType = "normal") {
    const dlurl = `http://localhost:8080/static/downloads/${filename}`;
    const posturl = "http://localhost:8080/upload";

    try {
        let result = await downloadAndPostFile(dlurl, posturl, uploadType);
        console.log("Result from server:", result);
    } catch (error) {
        console.error("Error during file download and posting:", error);
    }
}

async function downloadAndPostFile(fileUrl, postUrl, uploadType) {
    const filename = fileUrl.substring(fileUrl.lastIndexOf("/") + 1);
    const statusSpan = document.getElementById(`${filename}-status`);

    try {
        const blob = await downloadFileToBlob(fileUrl);

        if (uploadType === "normal") {
            await normalUpload(blob, filename, postUrl, statusSpan);
        } else if (uploadType === "secure") {
            await secureUpload(blob, filename, statusSpan);
        }

        console.log("File posted successfully");
    } catch (error) {
        if (statusSpan) {
            statusSpan.textContent = `POST not successful - ${error.message}`;
        }
        console.error("Error downloading or posting file:", error);
    }
}

async function downloadFileToBlob(fileUrl) {
    const response = await fetch(fileUrl);
    if (!response.ok) throw new Error(`Failed to download file. Status: ${response.status}`);
    return await response.blob();
}

async function normalUpload(file, filename, postUrl, statusSpan) {
    const formData = new FormData();
    formData.append("file", file, filename);

    try {
        const response = await fetch(postUrl, { method: "POST", body: formData });

        if (!response.ok) throw new Error(`POST error: ${response.status}`);

        if (statusSpan) {
            statusSpan.textContent = "POST successful";
        }

        return await response.json();
    } catch (error) {
        if (statusSpan) {
            statusSpan.textContent = `POST not successful - ${error.message}`;
        }
        console.error("Error during normal upload:", error);
    }
}

async function secureUpload(file, filename, statusSpan) {
    try {
        const response = await fetch("/generate-s3-url");
        if (!response.ok) throw new Error("Failed to retrieve S3 URL.");

        const data = await response.json();
        const uploadUrl = data.url;

        const uploadResponse = await fetch(uploadUrl, {
            method: "PUT",
            body: file,
            headers: { "Content-Type": file.type }
        });

        if (!uploadResponse.ok) throw new Error("Failed to upload file to S3.");

        if (statusSpan) {
            statusSpan.textContent = "File securely uploaded to S3!";
        }
    } catch (error) {
        if (statusSpan) {
            statusSpan.textContent = `Error during secure upload: ${error.message}`;
        }
        console.error("Error during secure upload:", error);
    }
}
