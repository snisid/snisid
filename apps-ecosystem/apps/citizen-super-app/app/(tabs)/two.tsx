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
import DocumentCard from '@/components/DocumentCard';
import type { Document } from '@/lib/api';

export default function DocumentsScreen() {
  const { api } = useAuth();
  const [documents, setDocuments] = useState<Document[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState('');

  const fetchDocuments = useCallback(async () => {
    try {
      setError('');
      const response = await api.getDocuments();
      setDocuments(response.documents);
    } catch (err: any) {
      setError(err?.response?.data?.message ?? err?.message ?? 'Failed to load documents');
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  }, [api]);

  useEffect(() => {
    fetchDocuments();
  }, [fetchDocuments]);

  const onRefresh = useCallback(() => {
    setRefreshing(true);
    fetchDocuments();
  }, [fetchDocuments]);

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#0033a0" />
        <Text style={styles.loadingText}>Loading documents...</Text>
      </View>
    );
  }

  if (error) {
    return (
      <View style={styles.center}>
        <Text style={styles.errorIcon}>⚠️</Text>
        <Text style={styles.errorText}>{error}</Text>
        <TouchableOpacity style={styles.retryButton} onPress={fetchDocuments}>
          <Text style={styles.retryText}>Retry</Text>
        </TouchableOpacity>
      </View>
    );
  }

  if (documents.length === 0) {
    return (
      <View style={styles.center}>
        <Text style={styles.emptyIcon}>📂</Text>
        <Text style={styles.emptyTitle}>No Documents Yet</Text>
        <Text style={styles.emptySubtitle}>
          Your digital documents will appear here once issued.
        </Text>
        <TouchableOpacity
          style={styles.retryButton}
          onRefresh={onRefresh}
          onPress={fetchDocuments}>
          <Text style={styles.retryText}>Refresh</Text>
        </TouchableOpacity>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <FlatList
        data={documents}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => <DocumentCard document={item} />}
        contentContainerStyle={styles.list}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
        ListHeaderComponent={
          <Text style={styles.header}>
            {documents.length} Document{documents.length !== 1 ? 's' : ''}
          </Text>
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
    gap: 8,
  },
  center: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 24,
    backgroundColor: '#F5F5F7',
  },
  header: {
    fontSize: 20,
    fontWeight: '700',
    color: '#000',
    paddingHorizontal: 20,
    marginBottom: 12,
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
    marginBottom: 24,
    lineHeight: 20,
  },
});
