document.addEventListener('DOMContentLoaded', () => {
    const generateBtn = document.getElementById('generate-btn');
    const qrCodePreview = document.getElementById('qr-code-preview');
    const downloadBtn = document.getElementById('download-btn');
    const textInput = document.getElementById('text');
    const logoInput = document.getElementById('logo');
    const bgColorInput = document.getElementById('bg-color');
    const dotColorInput = document.getElementById('dot-color');
    const eyeColorInput = document.getElementById('eye-color');
    const eyeShapeInput = document.getElementById('eye-shape');
    const dotShapeInput = document.getElementById('dot-shape');
    const paddingInput = document.getElementById('padding');

    generateBtn.addEventListener('click', () => {
        const text = textInput.value;
        if (!text) {
            alert('Vui lòng nhập nội dung hoặc URL.');
            return;
        }

        qrCodePreview.innerHTML = '<p class="text-gray-400">Đang tạo...</p>';

        const formData = new FormData();
        formData.append('text', text);
        formData.append('bgColor', bgColorInput.value);
        formData.append('dotColor', dotColorInput.value);
        formData.append('eyeColor', eyeColorInput.value);
        formData.append('eyeShape', eyeShapeInput.value);
        formData.append('dotShape', dotShapeInput.value);
        formData.append('padding', paddingInput.value);

        if (logoInput.files && logoInput.files[0]) {
            formData.append('logo', logoInput.files[0]);
        }

        const apiUrl = '/api/qr';

        fetch(apiUrl, {
            method: 'POST',
            body: formData,
        })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => { throw new Error(text) });
            }
            return response.blob();
        })
        .then(blob => {
            const imageUrl = URL.createObjectURL(blob);
            qrCodePreview.innerHTML = `<img src="${imageUrl}" alt="Generated QR Code" class="qr-image">`;
            downloadBtn.href = imageUrl;
            downloadBtn.classList.add('visible');
        })
        .catch(error => {
            console.error('Error generating QR code:', error);
            qrCodePreview.innerHTML = `<p class="text-red-500">Lỗi: ${error.message}</p>`;
            downloadBtn.classList.remove('visible');
        });
    });
});