import { useCallback, useEffect, useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ActivityIndicator,
  ScrollView,
  TouchableOpacity,
  RefreshControl,
} from 'react-native';
import { useRouter } from 'expo-router';

import { useAuth } from '@/lib/auth';
import IdentityCard from '@/components/IdentityCard';
import type { Identity } from '@/lib/api';

export default function IdentityDashboard() {
  const router = useRouter();
  const { api } = useAuth();
  const [identity, setIdentity] = useState<Identity | null>(null);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState('');

  const fetchIdentity = useCallback(async () => {
    try {
      setError('');
      const response = await api.getIdentity();
      setIdentity(response.identity);
    } catch (err: any) {
      setError(err?.response?.data?.message ?? err?.message ?? 'Failed to load identity');
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  }, [api]);

  useEffect(() => {
    fetchIdentity();
  }, [fetchIdentity]);

  const onRefresh = useCallback(() => {
    setRefreshing(true);
    fetchIdentity();
  }, [fetchIdentity]);

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#0033a0" />
        <Text style={styles.loadingText}>Loading your identity...</Text>
      </View>
    );
  }

  if (error) {
    return (
      <View style={styles.center}>
        <Text style={styles.errorIcon}>⚠️</Text>
        <Text style={styles.errorText}>{error}</Text>
        <TouchableOpacity style={styles.retryButton} onPress={fetchIdentity}>
          <Text style={styles.retryText}>Retry</Text>
        </TouchableOpacity>
      </View>
    );
  }

  return (
    <ScrollView
      style={styles.container}
      contentContainerStyle={styles.content}
      refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}>
      <Text style={styles.greeting}>Welcome back</Text>
      <Text style={styles.name}>{identity?.fullName}</Text>

      {identity && <IdentityCard identity={identity} />}

      <TouchableOpacity
        style={styles.verifyButton}
        onPress={() => router.push('/modal')}>
        <Text style={styles.verifyIcon}>📷</Text>
        <Text style={styles.verifyText}>Verify Identity</Text>
        <Text style={styles.verifyArrow}>→</Text>
      </TouchableOpacity>

      <View style={styles.quickActions}>
        <TouchableOpacity style={styles.actionCard}>
          <Text style={styles.actionIcon}>🪪</Text>
          <Text style={styles.actionLabel}>My Documents</Text>
        </TouchableOpacity>
        <TouchableOpacity style={styles.actionCard}>
          <Text style={styles.actionIcon}>🔒</Text>
          <Text style={styles.actionLabel}>Privacy Settings</Text>
        </TouchableOpacity>
        <TouchableOpacity style={styles.actionCard}>
          <Text style={styles.actionIcon}>📞</Text>
          <Text style={styles.actionLabel}>Support</Text>
        </TouchableOpacity>
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
    paddingBottom: 32,
    gap: 20,
  },
  center: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 24,
    backgroundColor: '#F5F5F7',
  },
  greeting: {
    fontSize: 16,
    color: '#666',
    paddingHorizontal: 20,
  },
  name: {
    fontSize: 28,
    fontWeight: '700',
    color: '#000',
    paddingHorizontal: 20,
    marginBottom: 8,
  },
  loadingText: {
    marginTop: 12,
    fontSize: 14,
    color: '#666',
  },
  errorIcon: {
    fontSize: 48,
    marginBottom: 12,
  },
  errorText: {
    fontSize: 14,
    color: '#991B1B',
    textAlign: 'center',
    marginBottom: 16,
  },
  retryButton: {
    backgroundColor: '#0033a0',
    paddingHorizontal: 24,
    paddingVertical: 12,
    borderRadius: 8,
  },
  retryText: {
    color: '#fff',
    fontWeight: '600',
  },
  verifyButton: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#0033a0',
    marginHorizontal: 16,
    padding: 16,
    borderRadius: 14,
  },
  verifyIcon: {
    fontSize: 24,
    marginRight: 12,
  },
  verifyText: {
    flex: 1,
    color: '#fff',
    fontSize: 17,
    fontWeight: '600',
  },
  verifyArrow: {
    color: '#fff',
    fontSize: 20,
  },
  quickActions: {
    flexDirection: 'row',
    paddingHorizontal: 16,
    gap: 12,
  },
  actionCard: {
    flex: 1,
    backgroundColor: '#fff',
    borderRadius: 14,
    padding: 16,
    alignItems: 'center',
    gap: 8,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.05,
    shadowRadius: 4,
    elevation: 2,
  },
  actionIcon: {
    fontSize: 28,
  },
  actionLabel: {
    fontSize: 12,
    fontWeight: '600',
    color: '#333',
    textAlign: 'center',
  },
});
