import { useCallback, useEffect, useState } from 'react';
import {
  View,
  Text,
  FlatList,
  StyleSheet,
  ActivityIndicator,
  TouchableOpacity,
  RefreshControl,
} from 'react-native';

import { useAuth } from '@/lib/auth';
import ActivityCard from '@/components/ActivityCard';
import type { NotificationItem } from '@/lib/api';

const TYPE_CONFIG: Record<
  string,
  { icon: string; status: 'success' | 'warning' | 'error' | 'info' }
> = {
  identity_update: { icon: '🪪', status: 'info' },
  verification: { icon: '✅', status: 'success' },
  document_expiry: { icon: '⏰', status: 'warning' },
  security_alert: { icon: '🔒', status: 'error' },
  system: { icon: '⚙️', status: 'info' },
};

export default function NotificationsScreen() {
  const { api } = useAuth();
  const [notifications, setNotifications] = useState<NotificationItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState('');
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [loadingMore, setLoadingMore] = useState(false);

  const fetchNotifications = useCallback(
    async (pageNum = 1, append = false) => {
      try {
        if (append) setLoadingMore(true);
        else setError('');
        const response = await api.getNotifications({ page: pageNum, limit: 20 });
        if (append) {
          setNotifications((prev) => [...prev, ...response.notifications]);
        } else {
          setNotifications(response.notifications);
        }
        setTotal(response.total);
        setPage(pageNum);
      } catch (err: any) {
        if (!append) {
          setError(
            err?.response?.data?.message ?? err?.message ?? 'Failed to load notifications'
          );
        }
      } finally {
        setLoading(false);
        setRefreshing(false);
        setLoadingMore(false);
      }
    },
    [api]
  );

  useEffect(() => {
    fetchNotifications();
  }, [fetchNotifications]);

  const onRefresh = useCallback(() => {
    setRefreshing(true);
    fetchNotifications(1);
  }, [fetchNotifications]);

  const onEndReached = useCallback(() => {
    if (!loadingMore && notifications.length < total) {
      fetchNotifications(page + 1, true);
    }
  }, [loadingMore, notifications.length, total, page, fetchNotifications]);

  const handleMarkRead = useCallback(
    async (id: string) => {
      try {
        await api.markNotificationRead(id);
        setNotifications((prev) =>
          prev.map((n) => (n.id === id ? { ...n, read: true } : n))
        );
      } catch {}
    },
    [api]
  );

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#0033a0" />
        <Text style={styles.loadingText}>Loading notifications...</Text>
      </View>
    );
  }

  if (error) {
    return (
      <View style={styles.center}>
        <Text style={styles.errorIcon}>⚠️</Text>
        <Text style={styles.errorText}>{error}</Text>
        <TouchableOpacity style={styles.retryButton} onPress={() => fetchNotifications()}>
          <Text style={styles.retryText}>Retry</Text>
        </TouchableOpacity>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <FlatList
        data={notifications}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => {
          const config = TYPE_CONFIG[item.type] ?? TYPE_CONFIG.system;
          return (
            <TouchableOpacity
              activeOpacity={0.7}
              onPress={() => !item.read && handleMarkRead(item.id)}>
              <View style={[!item.read && styles.unreadBorder]}>
                <ActivityCard
                  icon={config.icon}
                  title={item.title}
                  description={item.description}
                  timestamp={item.timestamp}
                  status={config.status}
                />
              </View>
            </TouchableOpacity>
          );
        }}
        contentContainerStyle={styles.list}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
        }
        onEndReached={onEndReached}
        onEndReachedThreshold={0.3}
        ListHeaderComponent={
          notifications.length > 0 ? (
            <View style={styles.headerRow}>
              <Text style={styles.header}>
                {total} Notification{total !== 1 ? 's' : ''}
              </Text>
            </View>
          ) : null
        }
        ListFooterComponent={
          loadingMore ? (
            <View style={styles.footerLoader}>
              <ActivityIndicator size="small" color="#0033a0" />
            </View>
          ) : null
        }
        ListEmptyComponent={
          <View style={styles.emptyContainer}>
            <Text style={styles.emptyIcon}>🔔</Text>
            <Text style={styles.emptyTitle}>No Notifications</Text>
            <Text style={styles.emptySubtitle}>
              You're all caught up! New activity will appear here.
            </Text>
          </View>
        }
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#F5F5F7',
  },
  list: {
    paddingVertical: 16,
  },
  center: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 24,
    backgroundColor: '#F5F5F7',
  },
  headerRow: {
    paddingHorizontal: 20,
    marginBottom: 8,
  },
  header: {
    fontSize: 20,
    fontWeight: '700',
    color: '#000',
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
  unreadBorder: {
    borderLeftWidth: 3,
    borderLeftColor: '#0033a0',
    marginLeft: 16,
    paddingLeft: 0,
  },
  footerLoader: {
    paddingVertical: 20,
    alignItems: 'center',
  },
  emptyContainer: {
    alignItems: 'center',
    paddingTop: 60,
    paddingHorizontal: 24,
  },
  emptyIcon: {
    fontSize: 64,
    marginBottom: 16,
  },
  emptyTitle: {
    fontSize: 20,
    fontWeight: '700',
    color: '#000',
    marginBottom: 8,
  },
  emptySubtitle: {
    fontSize: 14,
    color: '#666',
    textAlign: 'center',
    lineHeight: 20,
  },
});
