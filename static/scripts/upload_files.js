async function uploadTestFile() {
    try {
        // Step 1: Fetch the signed URL from the Go API
        const response = await fetch("https://your-api.com/get-signed-url");
        if (!response.ok) throw new Error("Failed to get signed URL");

        const { url } = await response.json(); // Extract signed URL

        // Step 2: Create a test file as a Blob (Example: a simple text file)
        const testFile = new Blob(["Test file content"], { type: "text/plain" });

        // Step 3: Upload the file to S3
        const putResponse = await fetch(url, {
            method: "PUT",
            body: testFile,
            headers: { "Content-Type": "text/plain" } // Adjust for your file type
        });

        if (putResponse.ok) {
            console.log("✅ File uploaded successfully!");
            alert("Upload Successful!");
        } else {
            throw new Error("S3 upload failed");
        }
    } catch (error) {
        console.error("Upload error:", error);
        alert("Upload Failed: " + error.message);
    }
}
