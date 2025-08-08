
// 1. 先从后端获取预授权上传 URL
const getPresignedUploadURL = async (fileName) => {
    // 调用后端 API 获取预授权 URL
    const response = await fetch('/api/get-presigned-upload-url', {
      method: 'POST',
      body: JSON.stringify({ fileName }),
      headers: { 'Content-Type': 'application/json' }
    });
    const data = await response.json();
    return data.uploadURL;
  };
  

  // 2. 在 NutUI Uploader 的 upload 函数中使用
  const upload = async (file) => {
    try {
      // 获取预授权 URL
      const uploadURL = await getPresignedUploadURL(file.name);
      
      // 使用 PUT 方法上传到预授权 URL
      const response = await fetch(uploadURL, {
        method: 'PUT',
        body: file,
        headers: {
          'Content-Type': file.type || 'application/octet-stream'
        }
      });
      
      if (response.ok) {
        return { url: uploadURL.split('?')[0] }; // 返回文件的访问 URL
      } else {
        throw new Error('Upload failed');
      }
    } catch (error) {
      throw error;
    }
  };