import { useState, useCallback } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ActivityIndicator,
  TouchableOpacity,
  ScrollView,
} from 'react-native';
import { StatusBar } from 'expo-status-bar';

import QRScanner from '@/components/QRScanner';
import { useAuth } from '@/lib/auth';
import IdentityCard from '@/components/IdentityCard';
import type { Identity } from '@/lib/api';

type ScanState = 'scanning' | 'verifying' | 'result' | 'error';

export default function ModalScreen() {
  const { api } = useAuth();
  const [scanState, setScanState] = useState<ScanState>('scanning');
  const [result, setResult] = useState<{
    verified: boolean;
    message: string;
    identity?: Identity;
  } | null>(null);
  const [errorMessage, setErrorMessage] = useState('');

  const handleScan = useCallback(
    async (data: string) => {
      setScanState('verifying');
      try {
        const response = await api.verifyIdentity({ qrData: data });
        setResult(response);
        setScanState('result');
      } catch (err: any) {
        setErrorMessage(
          err?.response?.data?.message ?? err?.message ?? 'Verification failed'
        );
        setScanState('error');
      }
    },
    [api]
  );

  const resetScan = () => {
    setScanState('scanning');
    setResult(null);
    setErrorMessage('');
  };

  return (
    <View style={styles.container}>
      <StatusBar style="light" />
      <View style={styles.header}>
        <Text style={styles.title}>Identity Verification</Text>
        <Text style={styles.subtitle}>
          {scanState === 'scanning'
            ? 'Scan the QR code on the identity document'
            : scanState === 'verifying'
              ? 'Verifying identity...'
              : 'Verification complete'}
        </Text>
      </View>

      {scanState === 'scanning' && (
        <View style={styles.scannerContainer}>
          <QRScanner onScan={handleScan} onError={setErrorMessage} />
        </View>
      )}

      {scanState === 'verifying' && (
        <View style={styles.centerContent}>
          <ActivityIndicator size="large" color="#0033a0" />
          <Text style={styles.verifyingText}>Verifying identity data...</Text>
        </View>
      )}

      {scanState === 'result' && result && (
        <ScrollView style={styles.resultContainer}>
          <View
            style={[
              styles.resultBanner,
              result.verified ? styles.resultSuccess : styles.resultFailure,
            ]}>
            <Text style={styles.resultIcon}>{result.verified ? '✅' : '❌'}</Text>
            <Text style={styles.resultTitle}>
              {result.verified ? 'Identity Verified' : 'Verification Failed'}
            </Text>
            <Text style={styles.resultMessage}>{result.message}</Text>
          </View>

          {result.identity && <IdentityCard identity={result.identity} />}

          <TouchableOpacity style={styles.scanAgainButton} onPress={resetScan}>
            <Text style={styles.scanAgainText}>Scan Another Code</Text>
          </TouchableOpacity>
        </ScrollView>
      )}

      {scanState === 'error' && (
        <View style={styles.centerContent}>
          <Text style={styles.errorIcon}>⚠️</Text>
          <Text style={styles.errorTitle}>Scan Error</Text>
          <Text style={styles.errorDetail}>{errorMessage}</Text>
          <TouchableOpacity style={styles.scanAgainButton} onPress={resetScan}>
            <Text style={styles.scanAgainText}>Try Again</Text>
          </TouchableOpacity>
        </View>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#000',
  },
  header: {
    paddingTop: 16,
    paddingHorizontal: 20,
    paddingBottom: 16,
    zIndex: 10,
  },
  title: {
    fontSize: 22,
    fontWeight: '700',
    color: '#fff',
  },
  subtitle: {
    fontSize: 14,
    color: '#aaa',
    marginTop: 4,
  },
  scannerContainer: {
    flex: 1,
    marginHorizontal: 12,
    borderRadius: 16,
    overflow: 'hidden',
  },
  centerContent: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 24,
  },
  verifyingText: {
    color: '#fff',
    fontSize: 16,
    marginTop: 16,
  },
  resultContainer: {
    flex: 1,
    paddingHorizontal: 12,
  },
  resultBanner: {
    borderRadius: 16,
    padding: 20,
    alignItems: 'center',
    marginBottom: 20,
  },
  resultSuccess: {
    backgroundColor: '#065F46',
  },
  resultFailure: {
    backgroundColor: '#7F1D1D',
  },
  resultIcon: {
    fontSize: 40,
    marginBottom: 12,
  },
  resultTitle: {
    fontSize: 20,
    fontWeight: '700',
    color: '#fff',
    marginBottom: 8,
  },
  resultMessage: {
    fontSize: 14,
    color: '#ddd',
    textAlign: 'center',
    lineHeight: 20,
  },
  scanAgainButton: {
    backgroundColor: '#0033a0',
    borderRadius: 12,
    paddingVertical: 14,
    alignItems: 'center',
    marginTop: 20,
    marginBottom: 40,
  },
  scanAgainText: {
    color: '#fff',
    fontSize: 17,
    fontWeight: '600',
  },
  errorIcon: {
    fontSize: 48,
    marginBottom: 12,
  },
  errorTitle: {
    fontSize: 20,
    fontWeight: '700',
    color: '#fff',
    marginBottom: 8,
  },
  errorDetail: {
    fontSize: 14,
    color: '#aaa',
    textAlign: 'center',
    marginBottom: 24,
  },
});
