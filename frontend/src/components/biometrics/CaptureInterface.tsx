import React, { useRef, useState, useCallback } from 'react';
import { Camera, RefreshCw, Check, AlertTriangle } from 'lucide-react';
import { logAuditAction } from '../../services/auditLogger';

interface CaptureInterfaceProps {
  onCapture: (imageBlob: Blob) => void;
  requiredQuality?: number;
}

export const CaptureInterface: React.FC<CaptureInterfaceProps> = ({ 
  onCapture, 
  requiredQuality = 0.8 
}) => {
  const videoRef = useRef<HTMLVideoElement>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  
  const [isStreaming, setIsStreaming] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [capturedImage, setCapturedImage] = useState<string | null>(null);
  const [qualityScore, setQualityScore] = useState<number | null>(null);

  const startCamera = async () => {
    try {
      setError(null);
      const stream = await navigator.mediaDevices.getUserMedia({ 
        video: { 
          width: { ideal: 1920 },
          height: { ideal: 1080 },
          facingMode: 'user'
        } 
      });
      
      if (videoRef.current) {
        videoRef.current.srcObject = stream;
        setIsStreaming(true);
      }
      
      // Log operator action
      logAuditAction('CAMERA_STARTED', { device: 'webcam' });
      
    } catch (err) {
      console.error('Error accessing camera:', err);
      setError("Impossible d'accéder à la caméra. Vérifiez vos autorisations.");
      logAuditAction('CAMERA_ACCESS_DENIED', { error: String(err) });
    }
  };

  const stopCamera = useCallback(() => {
    if (videoRef.current && videoRef.current.srcObject) {
      const stream = videoRef.current.srcObject as MediaStream;
      stream.getTracks().forEach(track => track.stop());
      videoRef.current.srcObject = null;
      setIsStreaming(false);
    }
  }, []);

  // Cleanup on unmount
  React.useEffect(() => {
    return () => stopCamera();
  }, [stopCamera]);

  const captureFrame = () => {
    if (!videoRef.current || !canvasRef.current) return;

    const video = videoRef.current;
    const canvas = canvasRef.current;
    
    canvas.width = video.videoWidth;
    canvas.height = video.videoHeight;
    
    const ctx = canvas.getContext('2d');
    if (!ctx) return;
    
    ctx.drawImage(video, 0, 0, canvas.width, canvas.height);
    
    // Simulate client-side quality check (e.g., face detection, blur detection)
    // In production, this would use a lightweight WASM model
    const simulatedQuality = Math.random() * 0.4 + 0.6; // Random between 0.6 and 1.0
    setQualityScore(simulatedQuality);
    
    const dataUrl = canvas.toDataURL('image/jpeg', 0.92);
    setCapturedImage(dataUrl);
    stopCamera();

    logAuditAction('BIOMETRIC_CAPTURED', { 
      type: 'face', 
      quality: simulatedQuality 
    });
  };

  const handleRetake = () => {
    setCapturedImage(null);
    setQualityScore(null);
    startCamera();
  };

  const handleConfirm = () => {
    if (!canvasRef.current || qualityScore === null) return;
    if (qualityScore < requiredQuality) {
      setError(`La qualité de l'image (${(qualityScore * 100).toFixed(1)}%) est inférieure au seuil requis (${requiredQuality * 100}%).`);
      return;
    }

    canvasRef.current.toBlob((blob) => {
      if (blob) {
        logAuditAction('BIOMETRIC_CONFIRMED', { quality: qualityScore });
        onCapture(blob);
      }
    }, 'image/jpeg', 0.92);
  };

  return (
    <div className="bg-slate-900 border border-slate-800 rounded-xl p-6 shadow-xl max-w-2xl w-full">
      <h3 className="text-xl font-bold text-white mb-4">Capture Biométrique</h3>
      
      {error && (
        <div className="mb-4 p-3 bg-red-500/10 border border-red-500/20 text-red-500 rounded-lg flex items-center gap-2 text-sm">
          <AlertTriangle size={16} />
          {error}
        </div>
      )}

      <div className="relative aspect-video bg-black rounded-lg overflow-hidden border-2 border-slate-700 mb-6 flex items-center justify-center">
        {!isStreaming && !capturedImage && (
          <div className="text-center">
            <Camera size={48} className="mx-auto text-slate-600 mb-2" />
            <p className="text-slate-400">Caméra désactivée</p>
          </div>
        )}
        
        <video 
          ref={videoRef} 
          autoPlay 
          playsInline 
          muted 
          className={`absolute inset-0 w-full h-full object-cover ${capturedImage ? 'hidden' : 'block'}`}
        />
        
        {capturedImage && (
          <img 
            src={capturedImage} 
            alt="Capture" 
            className="absolute inset-0 w-full h-full object-cover" 
          />
        )}
        
        {/* Simulated Face Overlay Guideline */}
        {isStreaming && (
          <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
            <div className="w-1/3 h-1/2 border-2 border-dashed border-blue-500/50 rounded-full opacity-50"></div>
          </div>
        )}

        {/* Quality Indicator */}
        {capturedImage && qualityScore !== null && (
          <div className={`absolute top-4 right-4 px-3 py-1.5 rounded-lg font-bold text-sm backdrop-blur-md border ${
            qualityScore >= requiredQuality 
              ? 'bg-green-500/20 text-green-400 border-green-500/30' 
              : 'bg-red-500/20 text-red-400 border-red-500/30'
          }`}>
            Qualité: {(qualityScore * 100).toFixed(1)}%
          </div>
        )}
      </div>

      <div className="flex gap-4 justify-center">
        {!isStreaming && !capturedImage ? (
          <button 
            onClick={startCamera}
            className="flex items-center gap-2 px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors"
          >
            <Camera size={20} />
            Démarrer la caméra
          </button>
        ) : isStreaming ? (
          <button 
            onClick={captureFrame}
            className="flex items-center gap-2 px-8 py-3 bg-white hover:bg-slate-200 text-black rounded-full font-bold transition-colors shadow-lg"
          >
            Capturer
          </button>
        ) : (
          <>
            <button 
              onClick={handleRetake}
              className="flex items-center gap-2 px-6 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg font-medium transition-colors"
            >
              <RefreshCw size={20} />
              Reprendre
            </button>
            <button 
              onClick={handleConfirm}
              disabled={qualityScore !== null && qualityScore < requiredQuality}
              className="flex items-center gap-2 px-6 py-2 bg-green-600 hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg font-medium transition-colors"
            >
              <Check size={20} />
              Confirmer
            </button>
          </>
        )}
      </div>

      {/* Hidden canvas for capturing frame */}
      <canvas ref={canvasRef} className="hidden" />
    </div>
  );
};
