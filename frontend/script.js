document.addEventListener('DOMContentLoaded', () => {
    const generateBtn = document.getElementById('generate-btn');
    const qrCodePreview = document.getElementById('qr-code-preview');
    const downloadBtn = document.getElementById('download-btn');
    const textInput = document.getElementById('text');
    const logoInput = document.getElementById('logo');
    const bgColorInput = document.getElementById('bg-color');
    const dotColorInput = document.getElementById('dot-color');

    generateBtn.addEventListener('click', () => {
        const text = textInput.value;
        if (!text) {
            alert('Please enter text or a URL.');
            return;
        }

        qrCodePreview.innerHTML = '<p class="text-gray-400">Generating...</p>';

        const formData = new FormData();
        formData.append('text', text);
        formData.append('bgColor', bgColorInput.value);
        formData.append('dotColor', dotColorInput.value);

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
            qrCodePreview.innerHTML = `<img src="${imageUrl}" alt="Generated QR Code" class="w-full h-full">`;
            downloadBtn.href = imageUrl;
            downloadBtn.classList.remove('hidden');
        })
        .catch(error => {
            console.error('Error generating QR code:', error);
            qrCodePreview.innerHTML = `<p class="text-red-500">Error: ${error.message}</p>`;
            downloadBtn.classList.add('hidden');
        });
    });
});