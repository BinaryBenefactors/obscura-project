import React, { useState, useRef } from 'react';
import '../styles/style.css';

const UploadPage = () => {
    const [selectedFile, setSelectedFile] = useState(null);
    const [fileUrl, setFileUrl] = useState(null);
    const [isMediaShown, setIsMediaShown] = useState(false);
    const [progress, setProgress] = useState(0);
    const [statusMessage, setStatusMessage] = useState('');
    const [statusType, setStatusType] = useState('');
    
    const fileInputRef = useRef();
    const uploadAreaRef = useRef();
    const mediaPlayerContainerRef = useRef();

    const handleDragOver = (e) => {
        e.preventDefault();
        uploadAreaRef.current.classList.add('highlight');
    };

    const handleDragLeave = () => {
        uploadAreaRef.current.classList.remove('highlight');
    };

    const handleDrop = (e) => {
        e.preventDefault();
        uploadAreaRef.current.classList.remove('highlight');
        if (e.dataTransfer.files.length) {
            handleFileSelect(e.dataTransfer.files[0]);
        }
    };

    const handleFileSelect = (file) => {
        if (!file.type.match('image.*') && !file.type.match('video.*')) {
            showStatus('Пожалуйста, выберите изображение или видео', 'error');
            return;
        }
        setSelectedFile(file);
        setFileUrl(null);
        setIsMediaShown(false);
        setStatusMessage('');
        setProgress(0);
    };

    const uploadFile = async () => {
        if (!selectedFile) return;
        const formData = new FormData();
        formData.append('mediaFile', selectedFile);

        try {
            setProgress(0);
            setStatusMessage('');
            setStatusType('');
            const response = await fetch('/upload', {
                method: 'POST',
                body: formData,
            });

            if (response.status === 201) {
                const data = await response.json();
                setFileUrl(data.fileUrl);
                showStatus('Файл успешно загружен! Нажмите "Просмотреть"', 'success');
            } else {
                const error = await response.json();
                showStatus(error.error || 'Ошибка при загрузке файла', 'error');
            }
        } catch (err) {
            showStatus('Ошибка соединения', 'error');
        }
    };

    const showStatus = (message, type) => {
        setStatusMessage(message);
        setStatusType(type);
    };

    const toggleMediaPlayer = () => {
        if (!fileUrl) return;
        if (isMediaShown) {
            setIsMediaShown(false);
        } else {
            setIsMediaShown(true);
        }
    };

    const formatFileSize = (bytes) => {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    };

    return (
        <div className="container">
            <h1>Загрузка фото и видео</h1>

            <div
                className="upload-area"
                ref={uploadAreaRef}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
            >
                <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                    <polyline points="17 8 12 3 7 8"></polyline>
                    <line x1="12" y1="3" x2="12" y2="15"></line>
                </svg>
                <h3>Перетащите файлы сюда или</h3>
                <input
                    type="file"
                    ref={fileInputRef}
                    className="file-input"
                    accept="image/*,video/*"
                    onChange={(e) => {
                        if (e.target.files.length) {
                            handleFileSelect(e.target.files[0]);
                        }
                    }}
                />
                <button className="btn" onClick={() => fileInputRef.current.click()}>Выбрать файл</button>
                <div className="file-info">
                    {selectedFile
                        ? `${selectedFile.name} (${formatFileSize(selectedFile.size)})`
                        : 'Файл не выбран'}
                </div>
            </div>

            {selectedFile && (
                <>
                    <div className="progress-container">
                        <div className="progress-bar">
                            <div className="progress" style={{ width: `${progress}%` }}></div>
                        </div>
                        <div className="progress-text">{progress}%</div>
                    </div>

                    <div className="controls">
                        <button className="btn" onClick={uploadFile}>Загрузить</button>
                    </div>
                </>
            )}

            {statusMessage && (
                <div className={`status ${statusType}`}>{statusMessage}</div>
            )}

            {fileUrl && (
                <div className="preview-container">
                    {isMediaShown && (
                        <div ref={mediaPlayerContainerRef}>
                            {selectedFile.type.startsWith('image') ? (
                                <img src={fileUrl} alt="Загруженное изображение" className="media-player" />
                            ) : (
                                <video src={fileUrl} controls autoPlay className="media-player"></video>
                            )}
                        </div>
                    )}
                    <button
                        className={`btn ${isMediaShown ? 'btn-hide' : 'btn-view'}`}
                        onClick={toggleMediaPlayer}
                    >
                        {isMediaShown ? 'Спрятать' : 'Просмотреть'}
                    </button>
                </div>
            )}
        </div>
    );
};

export default UploadPage;
