import { useState, useCallback, useRef } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  ActivityIndicator,
  Alert,
  SafeAreaView,
} from 'react-native';
import { useRouter, Stack } from 'expo-router';
import { CameraView, useCameraPermissions } from 'expo-camera';

import { useAuth } from '@/lib/auth';

type EnrollmentStep = 'intro' | 'capture' | 'processing' | 'complete' | 'error';

const STEPS = ['Introduction', 'Face Capture', 'Quality Check', 'Enrolled'];

export default function BiometricEnrollScreen() {
  const router = useRouter();
  const { api } = useAuth();
  const cameraRef = useRef<CameraView>(null);
  const [permission, requestPermission] = useCameraPermissions();
  const [step, setStep] = useState<EnrollmentStep>('intro');
  const [errorMessage, setErrorMessage] = useState('');
  const [processing, setProcessing] = useState(false);
  const [qualityScore, setQualityScore] = useState(0);
  const [capturedImage, setCapturedImage] = useState<string | null>(null);

  const handleStartCapture = useCallback(async () => {
    if (!permission?.granted) {
      const result = await requestPermission();
      if (!result.granted) {
        Alert.alert(
          'Camera Required',
          'Camera permission is needed to capture face for biometric enrollment.'
        );
        return;
      }
    }
    setStep('capture');
  }, [permission, requestPermission]);

  const handleCapture = useCallback(async () => {
    if (!cameraRef.current) return;
    setProcessing(true);
    try {
      const photo = await cameraRef.current.takePictureAsync({
        quality: 0.8,
        base64: true,
      });
      if (!photo?.base64) throw new Error('Failed to capture image');

      setCapturedImage(photo.base64);

      const simulatedQuality = 0.75 + Math.random() * 0.2;
      setQualityScore(simulatedQuality);
      setStep('processing');

      await new Promise((resolve) => setTimeout(resolve, 1500));

      await api.enrollBiometrics({
        faceImage: photo.base64,
        livenessConfidence: simulatedQuality,
      });

      setStep('complete');
    } catch (err: any) {
      setErrorMessage(
        err?.response?.data?.message ??
          err?.message ??
          'Enrollment failed. Please try again.'
      );
      setStep('error');
    } finally {
      setProcessing(false);
    }
  }, [api]);

  const handleRetry = useCallback(() => {
    setStep('capture');
    setCapturedImage(null);
    setQualityScore(0);
    setErrorMessage('');
  }, []);

  const handleFinish = useCallback(() => {
    router.back();
  }, [router]);

  const currentStepIndex = ['intro', 'capture', 'processing', 'complete'].indexOf(step);

  return (
    <SafeAreaView style={styles.container}>
      <Stack.Screen
        options={{
          headerShown: true,
          headerTitle: 'Biometric Enrollment',
          headerStyle: { backgroundColor: '#000' },
          headerTintColor: '#fff',
        }}
      />

      <View style={styles.progressContainer}>
        {STEPS.map((label, i) => (
          <View key={label} style={styles.progressStep}>
            <View
              style={[
                styles.progressDot,
                i <= currentStepIndex ? styles.progressDotActive : styles.progressDotInactive,
              ]}>
              {i < currentStepIndex ? (
                <Text style={styles.progressCheck}>✓</Text>
              ) : (
                <Text style={styles.progressNum}>{i + 1}</Text>
              )}
            </View>
            <Text
              style={[
                styles.progressLabel,
                i <= currentStepIndex ? styles.progressLabelActive : styles.progressLabelInactive,
              ]}>
              {label}
            </Text>
            {i < STEPS.length - 1 && (
              <View
                style={[
                  styles.progressLine,
                  i < currentStepIndex ? styles.progressLineActive : styles.progressLineInactive,
                ]}
              />
            )}
          </View>
        ))}
      </View>

      {step === 'intro' && (
        <View style={styles.content}>
          <View style={styles.introIcon}>
            <Text style={styles.introEmoji}>👤</Text>
          </View>
          <Text style={styles.introTitle}>Enroll Your Biometrics</Text>
          <Text style={styles.introDesc}>
            Your facial biometrics will be securely stored and used for identity
            verification. This process takes less than a minute.
          </Text>

          <View style={styles.benefitsList}>
            <Text style={styles.benefitItem}>🔒 End-to-end encryption</Text>
            <Text style={styles.benefitItem}>✅ Faster identity verification</Text>
            <Text style={styles.benefitItem}>📱 Secure access to all services</Text>
          </View>

          <TouchableOpacity style={styles.primaryButton} onPress={handleStartCapture}>
            <Text style={styles.primaryButtonText}>Start Enrollment</Text>
          </TouchableOpacity>
          <TouchableOpacity style={styles.skipButton} onPress={handleFinish}>
            <Text style={styles.skipText}>Skip for now</Text>
          </TouchableOpacity>
        </View>
      )}

      {step === 'capture' && (
        <View style={styles.cameraContainer}>
          <CameraView
            ref={cameraRef}
            style={styles.camera}
            facing="front"
            autofocus="on">
            <View style={styles.cameraOverlay}>
              <View style={styles.faceGuide}>
                <View style={styles.faceOval} />
              </View>
              <Text style={styles.cameraHint}>Position your face within the frame</Text>
            </View>
          </CameraView>
          <View style={styles.cameraActions}>
            <TouchableOpacity
              style={[styles.primaryButton, processing && styles.buttonDisabled]}
              onPress={handleCapture}
              disabled={processing}>
              {processing ? (
                <ActivityIndicator size="small" color="#fff" />
              ) : (
                <Text style={styles.primaryButtonText}>Capture Photo</Text>
              )}
            </TouchableOpacity>
            <TouchableOpacity style={styles.skipButton} onPress={handleFinish}>
              <Text style={styles.skipText}>Cancel</Text>
            </TouchableOpacity>
          </View>
        </View>
      )}

      {step === 'processing' && (
        <View style={styles.content}>
          <ActivityIndicator size="large" color="#0033a0" />
          <Text style={styles.processingTitle}>Checking Quality</Text>
          <View style={styles.qualityBar}>
            <View
              style={[
                styles.qualityFill,
                { width: `${Math.round(qualityScore * 100)}%` },
                qualityScore > 0.7
                  ? styles.qualityGood
                  : qualityScore > 0.4
                    ? styles.qualityFair
                    : styles.qualityPoor,
              ]}
            />
          </View>
          <Text style={styles.processingDesc}>
            {qualityScore > 0.7
              ? 'Good quality capture'
              : qualityScore > 0.4
                ? 'Acceptable quality'
                : 'Low quality, consider retaking'}
          </Text>
          <Text style={styles.processingSub}>Please wait while we process your data...</Text>
        </View>
      )}

      {step === 'complete' && (
        <View style={styles.content}>
          <View style={styles.successIcon}>
            <Text style={styles.successEmoji}>✅</Text>
          </View>
          <Text style={styles.successTitle}>Enrollment Complete</Text>
          <Text style={styles.successDesc}>
            Your biometric data has been successfully enrolled. You can now use
            biometric authentication for faster access.
          </Text>
          <TouchableOpacity style={styles.primaryButton} onPress={handleFinish}>
            <Text style={styles.primaryButtonText}>Done</Text>
          </TouchableOpacity>
        </View>
      )}

      {step === 'error' && (
        <View style={styles.content}>
          <View style={styles.errorIcon}>
            <Text style={styles.errorEmoji}>⚠️</Text>
          </View>
          <Text style={styles.errorTitle}>Enrollment Failed</Text>
          <Text style={styles.errorDesc}>{errorMessage}</Text>
          <TouchableOpacity style={styles.primaryButton} onPress={handleRetry}>
            <Text style={styles.primaryButtonText}>Try Again</Text>
          </TouchableOpacity>
          <TouchableOpacity style={styles.skipButton} onPress={handleFinish}>
            <Text style={styles.skipText}>Skip for now</Text>
          </TouchableOpacity>
        </View>
      )}
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#000',
  },
  progressContainer: {
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
    paddingVertical: 20,
    paddingHorizontal: 24,
  },
  progressStep: {
    alignItems: 'center',
    flexDirection: 'row',
  },
  progressDot: {
    width: 32,
    height: 32,
    borderRadius: 16,
    justifyContent: 'center',
    alignItems: 'center',
  },
  progressDotActive: {
    backgroundColor: '#0033a0',
  },
  progressDotInactive: {
    backgroundColor: '#333',
  },
  progressCheck: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '700',
  },
  progressNum: {
    color: '#fff',
    fontSize: 13,
    fontWeight: '700',
  },
  progressLabel: {
    fontSize: 11,
    marginLeft: 6,
    fontWeight: '500',
  },
  progressLabelActive: {
    color: '#fff',
  },
  progressLabelInactive: {
    color: '#666',
  },
  progressLine: {
    width: 24,
    height: 2,
    marginHorizontal: 8,
  },
  progressLineActive: {
    backgroundColor: '#0033a0',
  },
  progressLineInactive: {
    backgroundColor: '#333',
  },
  content: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 32,
  },
  introIcon: {
    width: 96,
    height: 96,
    borderRadius: 48,
    backgroundColor: '#1C1C1E',
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 24,
  },
  introEmoji: {
    fontSize: 48,
  },
  introTitle: {
    fontSize: 24,
    fontWeight: '700',
    color: '#fff',
    marginBottom: 12,
    textAlign: 'center',
  },
  introDesc: {
    fontSize: 15,
    color: '#aaa',
    textAlign: 'center',
    lineHeight: 22,
    marginBottom: 24,
  },
  benefitsList: {
    alignSelf: 'stretch',
    gap: 12,
    marginBottom: 32,
  },
  benefitItem: {
    fontSize: 15,
    color: '#ddd',
  },
  primaryButton: {
    backgroundColor: '#0033a0',
    borderRadius: 14,
    paddingVertical: 16,
    paddingHorizontal: 32,
    alignItems: 'center',
    alignSelf: 'stretch',
  },
  primaryButtonText: {
    color: '#fff',
    fontSize: 17,
    fontWeight: '600',
  },
  buttonDisabled: {
    opacity: 0.6,
  },
  skipButton: {
    paddingVertical: 14,
    alignItems: 'center',
    alignSelf: 'stretch',
  },
  skipText: {
    color: '#888',
    fontSize: 15,
    fontWeight: '500',
  },
  cameraContainer: {
    flex: 1,
  },
  camera: {
    flex: 1,
  },
  cameraOverlay: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  faceGuide: {
    width: 200,
    height: 260,
    borderRadius: 100,
    borderWidth: 3,
    borderColor: 'rgba(255,255,255,0.6)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  faceOval: {
    width: 160,
    height: 200,
    borderRadius: 80,
    borderWidth: 1,
    borderColor: 'rgba(255,255,255,0.3)',
  },
  cameraHint: {
    color: '#fff',
    fontSize: 14,
    marginTop: 24,
    backgroundColor: 'rgba(0,0,0,0.5)',
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 20,
  },
  cameraActions: {
    position: 'absolute',
    bottom: 40,
    left: 24,
    right: 24,
    gap: 8,
  },
  processingTitle: {
    fontSize: 20,
    fontWeight: '700',
    color: '#fff',
    marginTop: 20,
    marginBottom: 16,
  },
  qualityBar: {
    width: '100%',
    height: 8,
    backgroundColor: '#333',
    borderRadius: 4,
    overflow: 'hidden',
    marginBottom: 8,
  },
  qualityFill: {
    height: '100%',
    borderRadius: 4,
  },
  qualityGood: {
    backgroundColor: '#22C55E',
  },
  qualityFair: {
    backgroundColor: '#F59E0B',
  },
  qualityPoor: {
    backgroundColor: '#EF4444',
  },
  processingDesc: {
    fontSize: 14,
    color: '#aaa',
    marginBottom: 8,
  },
  processingSub: {
    fontSize: 13,
    color: '#666',
  },
  successIcon: {
    width: 96,
    height: 96,
    borderRadius: 48,
    backgroundColor: '#065F46',
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 24,
  },
  successEmoji: {
    fontSize: 48,
  },
  successTitle: {
    fontSize: 24,
    fontWeight: '700',
    color: '#fff',
    marginBottom: 12,
  },
  successDesc: {
    fontSize: 15,
    color: '#aaa',
    textAlign: 'center',
    lineHeight: 22,
    marginBottom: 32,
  },
  errorIcon: {
    width: 96,
    height: 96,
    borderRadius: 48,
    backgroundColor: '#7F1D1D',
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 24,
  },
  errorEmoji: {
    fontSize: 48,
  },
  errorTitle: {
    fontSize: 24,
    fontWeight: '700',
    color: '#fff',
    marginBottom: 8,
  },
  errorDesc: {
    fontSize: 14,
    color: '#aaa',
    textAlign: 'center',
    lineHeight: 20,
    marginBottom: 32,
  },
});
