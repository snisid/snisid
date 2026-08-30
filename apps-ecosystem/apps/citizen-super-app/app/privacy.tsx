import { useState, useCallback, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  Switch,
  ActivityIndicator,
  Alert,
  RefreshControl,
} from 'react-native';
import * as LocalAuthentication from 'expo-local-authentication';
import * as SecureStore from 'expo-secure-store';
import { Platform } from 'react-native';

import { useAuth } from '@/lib/auth';
import type { AccessLog, ActiveSession } from '@/lib/api';

function SessionRow({
  session,
  onTerminate,
}: {
  session: ActiveSession;
  onTerminate: (id: string) => void;
}) {
  return (
    <View style={styles.sessionRow}>
      <View style={styles.sessionInfo}>
        <Text style={styles.sessionDevice}>{session.deviceName}</Text>
        <Text style={styles.sessionDetail}>
          {session.deviceType} · Last active:{' '}
          {new Date(session.lastActive).toLocaleDateString()}
        </Text>
        <Text style={styles.sessionDetail}>
          Logged in: {new Date(session.loggedInAt).toLocaleString()}
        </Text>
      </View>
      {session.current ? (
        <View style={styles.currentBadge}>
          <Text style={styles.currentBadgeText}>Current</Text>
        </View>
      ) : (
        <TouchableOpacity
          style={styles.terminateButton}
          onPress={() => onTerminate(session.id)}>
          <Text style={styles.terminateText}>Revoke</Text>
        </TouchableOpacity>
      )}
    </View>
  );
}

export default function PrivacyScreen() {
  const { api } = useAuth();
  const [logs, setLogs] = useState<AccessLog[]>([]);
  const [sessions, setSessions] = useState<ActiveSession[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState('');
  const [biometricEnabled, setBiometricEnabled] = useState(false);
  const [biometricAvailable, setBiometricAvailable] = useState(false);

  const checkBiometricStatus = useCallback(async () => {
    const hasHardware = await LocalAuthentication.hasHardwareAsync();
    const enrolled = await LocalAuthentication.isEnrolledAsync();
    setBiometricAvailable(hasHardware && enrolled);
    if (Platform.OS === 'web') {
      setBiometricEnabled(localStorage.getItem('biometric_enabled') === 'true');
    } else {
      const val = await SecureStore.getItemAsync('biometric_enabled');
      setBiometricEnabled(val === 'true');
    }
  }, []);

  const fetchData = useCallback(async () => {
    try {
      setError('');
      const [logRes, sessionRes] = await Promise.all([
        api.getAccessLogs({ limit: 20 }),
        api.getActiveSessions(),
      ]);
      setLogs(logRes.logs);
      setSessions(sessionRes.sessions);
    } catch (err: any) {
      setError(
        err?.response?.data?.message ?? err?.message ?? 'Failed to load privacy data'
      );
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  }, [api]);

  useEffect(() => {
    checkBiometricStatus();
    fetchData();
  }, [checkBiometricStatus, fetchData]);

  const onRefresh = useCallback(() => {
    setRefreshing(true);
    checkBiometricStatus();
    fetchData();
  }, [checkBiometricStatus, fetchData]);

  const handleBiometricToggle = useCallback(
    async (value: boolean) => {
      if (value && biometricAvailable) {
        const result = await LocalAuthentication.authenticateAsync({
          promptMessage: 'Authenticate to enable biometric login',
          fallbackLabel: 'Use password',
        });
        if (!result.success) {
          Alert.alert('Authentication Failed', 'Could not verify your identity.');
          return;
        }
      }
      if (Platform.OS === 'web') {
        localStorage.setItem('biometric_enabled', value ? 'true' : 'false');
      } else {
        await SecureStore.setItemAsync('biometric_enabled', value ? 'true' : 'false');
      }
      setBiometricEnabled(value);
    },
    [biometricAvailable]
  );

  const handleChangePin = useCallback(() => {
    Alert.alert(
      'Change PIN',
      'A verification link will be sent to your registered email address.',
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Send Link',
          onPress: () => Alert.alert('Email Sent', 'Check your inbox for the PIN reset link.'),
        },
      ]
    );
  }, []);

  const handleTerminateSession = useCallback(
    async (sessionId: string) => {
      try {
        await api.terminateSession(sessionId);
        setSessions((prev) => prev.filter((s) => s.id !== sessionId));
        Alert.alert('Session Revoked', 'The session has been terminated.');
      } catch (err: any) {
        Alert.alert(
          'Error',
          err?.response?.data?.message ?? err?.message ?? 'Failed to terminate session'
        );
      }
    },
    [api]
  );

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#0033a0" />
        <Text style={styles.loadingText}>Loading privacy settings...</Text>
      </View>
    );
  }

  return (
    <ScrollView
      style={styles.container}
      contentContainerStyle={styles.content}
      refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}>
      <Text style={styles.pageTitle}>Privacy & Security</Text>

      {error && (
        <View style={styles.errorBanner}>
          <Text style={styles.errorBannerText}>{error}</Text>
          <TouchableOpacity onPress={fetchData}>
            <Text style={styles.errorRetry}>Retry</Text>
          </TouchableOpacity>
        </View>
      )}

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Security Settings</Text>

        <View style={styles.settingCard}>
          <View style={styles.settingRow}>
            <View style={styles.settingInfo}>
              <Text style={styles.settingLabel}>Biometric Login</Text>
              <Text style={styles.settingDesc}>
                {biometricAvailable
                  ? 'Use fingerprint or face to sign in'
                  : 'Biometrics not available on this device'}
              </Text>
            </View>
            <Switch
              value={biometricEnabled}
              onValueChange={handleBiometricToggle}
              disabled={!biometricAvailable}
              trackColor={{ false: '#E5E5EA', true: '#0033a080' }}
              thumbColor={biometricEnabled ? '#0033a0' : '#f4f3f4'}
            />
          </View>

          <View style={styles.divider} />

          <TouchableOpacity style={styles.settingRow} onPress={handleChangePin}>
            <View style={styles.settingInfo}>
              <Text style={styles.settingLabel}>Change PIN</Text>
              <Text style={styles.settingDesc}>Update your account security PIN</Text>
            </View>
            <Text style={styles.arrow}>→</Text>
          </TouchableOpacity>
        </View>
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Active Sessions ({sessions.length})</Text>
        {sessions.length === 0 ? (
          <Text style={styles.emptyText}>No active sessions found.</Text>
        ) : (
          <View style={styles.sessionCard}>
            {sessions.map((s) => (
              <SessionRow
                key={s.id}
                session={s}
                onTerminate={handleTerminateSession}
              />
            ))}
          </View>
        )}
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Recent Data Access ({logs.length})</Text>
        {logs.length === 0 ? (
          <Text style={styles.emptyText}>No access logs recorded yet.</Text>
        ) : (
          logs.map((log) => (
            <View key={log.id} style={styles.logRow}>
              <View style={styles.logIcon}>
                <Text style={styles.logIconText}>🔍</Text>
              </View>
              <View style={styles.logInfo}>
                <Text style={styles.logApp}>{log.application}</Text>
                <Text style={styles.logAction}>{log.action}</Text>
                <Text style={styles.logMeta}>
                  {log.ipAddress}
                  {log.location ? ` · ${log.location}` : ''}
                </Text>
              </View>
              <Text style={styles.logTime}>
                {new Date(log.timestamp).toLocaleDateString()}
              </Text>
            </View>
          ))
        )}
      </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#F5F5F7',
  },
  content: {
    paddingTop: 24,
    paddingBottom: 40,
    gap: 24,
  },
  center: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 24,
    backgroundColor: '#F5F5F7',
  },
  pageTitle: {
    fontSize: 28,
    fontWeight: '700',
    color: '#000',
    paddingHorizontal: 20,
    marginBottom: 4,
  },
  loadingText: {
    marginTop: 12,
    fontSize: 14,
    color: '#666',
  },
  errorBanner: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#FEE2E2',
    marginHorizontal: 16,
    padding: 12,
    borderRadius: 10,
    gap: 8,
  },
  errorBannerText: {
    flex: 1,
    fontSize: 13,
    color: '#991B1B',
  },
  errorRetry: {
    fontSize: 13,
    fontWeight: '700',
    color: '#0033a0',
  },
  section: {
    paddingHorizontal: 16,
  },
  sectionTitle: {
    fontSize: 13,
    fontWeight: '600',
    color: '#888',
    textTransform: 'uppercase',
    letterSpacing: 0.5,
    marginBottom: 8,
    paddingHorizontal: 4,
  },
  settingCard: {
    backgroundColor: '#fff',
    borderRadius: 14,
    overflow: 'hidden',
  },
  settingRow: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 16,
  },
  settingInfo: {
    flex: 1,
    marginRight: 12,
  },
  settingLabel: {
    fontSize: 16,
    fontWeight: '600',
    color: '#000',
  },
  settingDesc: {
    fontSize: 12,
    color: '#888',
    marginTop: 2,
  },
  divider: {
    height: StyleSheet.hairlineWidth,
    backgroundColor: '#E5E5EA',
    marginHorizontal: 16,
  },
  arrow: {
    fontSize: 18,
    color: '#999',
  },
  sessionCard: {
    backgroundColor: '#fff',
    borderRadius: 14,
    overflow: 'hidden',
  },
  sessionRow: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 14,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: '#E5E5EA',
  },
  sessionInfo: {
    flex: 1,
    marginRight: 8,
  },
  sessionDevice: {
    fontSize: 15,
    fontWeight: '600',
    color: '#000',
  },
  sessionDetail: {
    fontSize: 11,
    color: '#888',
    marginTop: 2,
  },
  currentBadge: {
    backgroundColor: '#DCFCE7',
    paddingHorizontal: 10,
    paddingVertical: 4,
    borderRadius: 8,
  },
  currentBadgeText: {
    fontSize: 11,
    fontWeight: '700',
    color: '#166534',
  },
  terminateButton: {
    backgroundColor: '#FEE2E2',
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 8,
  },
  terminateText: {
    fontSize: 12,
    fontWeight: '700',
    color: '#991B1B',
  },
  logRow: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 12,
    marginBottom: 4,
  },
  logIcon: {
    width: 36,
    height: 36,
    borderRadius: 10,
    backgroundColor: '#F0F0F5',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 10,
  },
  logIconText: {
    fontSize: 16,
  },
  logInfo: {
    flex: 1,
    marginRight: 8,
  },
  logApp: {
    fontSize: 14,
    fontWeight: '600',
    color: '#000',
  },
  logAction: {
    fontSize: 12,
    color: '#555',
    marginTop: 1,
  },
  logMeta: {
    fontSize: 10,
    color: '#999',
    marginTop: 1,
  },
  logTime: {
    fontSize: 10,
    color: '#999',
  },
  emptyText: {
    fontSize: 14,
    color: '#999',
    textAlign: 'center',
    paddingVertical: 16,
  },
});
