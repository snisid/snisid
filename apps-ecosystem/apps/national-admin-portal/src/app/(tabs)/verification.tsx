import { useState, useCallback } from 'react';
import { FlatList, Pressable, RefreshControl, StyleSheet, View, Modal } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { SymbolView } from 'expo-symbols';

import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';

import { BottomTabInset, MaxContentWidth, Spacing } from '@/constants/theme';
import { useTheme } from '@/hooks/use-theme';

type VerificationRequest = {
  id: string;
  name: string;
  nnu: string;
  method: 'biometric' | 'document' | 'pin';
  submittedAt: string;
  agency: string;
};

const MOCK_REQUESTS: VerificationRequest[] = [
  { id: '1', name: 'Abubakar Sani', nnu: 'SN-88472', method: 'biometric', submittedAt: '5 min ago', agency: 'NIA HQ' },
  { id: '2', name: 'Chioma Okafor', nnu: 'SN-99103', method: 'document', submittedAt: '12 min ago', agency: 'Lagos State' },
  { id: '3', name: 'Emeka Nwosu', nnu: 'SN-77231', method: 'pin', submittedAt: '28 min ago', agency: 'FCT Admin' },
  { id: '4', name: 'Fatima Usman', nnu: 'SN-44567', method: 'biometric', submittedAt: '1 hour ago', agency: 'Kano State' },
  { id: '5', name: 'Gabriel Okonkwo', nnu: 'SN-33891', method: 'document', submittedAt: '2 hours ago', agency: 'NIA HQ' },
  { id: '6', name: 'Hauwa Ibrahim', nnu: 'SN-22904', method: 'biometric', submittedAt: '3 hours ago', agency: 'Rivers State' },
];

const METHOD_LABELS: Record<string, string> = { biometric: 'Biometric', document: 'Document', pin: 'PIN' };

export default function VerificationScreen() {
  const theme = useTheme();
  const [refreshing, setRefreshing] = useState(false);
  const [requests, setRequests] = useState(MOCK_REQUESTS);
  const [confirmModal, setConfirmModal] = useState<{ visible: boolean; id: string; action: 'approve' | 'reject' }>({ visible: false, id: '', action: 'approve' });

  const onRefresh = useCallback(() => {
    setRefreshing(true);
    setTimeout(() => setRefreshing(false), 1500);
  }, []);

  const handleAction = (id: string, action: 'approve' | 'reject') => {
    setConfirmModal({ visible: true, id, action });
  };

  const confirmAction = () => {
    setRequests(prev => prev.filter(r => r.id !== confirmModal.id));
    setConfirmModal({ visible: false, id: '', action: 'approve' });
  };

  return (
    <SafeAreaView style={[styles.safeArea, { backgroundColor: theme.background }]}>
      <FlatList
        data={requests}
        keyExtractor={(item) => item.id}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
        contentContainerStyle={styles.listContent}
        ListHeaderComponent={
          <ThemedView>
            <View style={styles.header}>
              <ThemedText type="subtitle" style={{ fontSize: 26, lineHeight: 32 }}>Verification Queue</ThemedText>
              <ThemedText themeColor="textSecondary" style={{ fontSize: 13 }}>{requests.length} pending</ThemedText>
            </View>
          </ThemedView>
        }
        renderItem={({ item }) => (
          <ThemedView type="backgroundElement" style={styles.card}>
            <View style={styles.cardRow}>
              <View style={[styles.avatar, { backgroundColor: theme.backgroundSelected }]}>
                <ThemedText style={{ fontWeight: '700', fontSize: 16 }}>{item.name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)}</ThemedText>
              </View>
              <View style={styles.cardInfo}>
                <ThemedText style={{ fontWeight: '600', fontSize: 15 }}>{item.name}</ThemedText>
                <ThemedText themeColor="textSecondary" style={{ fontSize: 13 }}>NNU: {item.nnu}</ThemedText>
                <ThemedText themeColor="textSecondary" style={{ fontSize: 12 }}>{item.agency} · {METHOD_LABELS[item.method]}</ThemedText>
              </View>
              <ThemedText themeColor="textSecondary" style={{ fontSize: 11 }}>{item.submittedAt}</ThemedText>
            </View>
            <View style={styles.actionRow}>
              <Pressable onPress={() => handleAction(item.id, 'reject')} style={[styles.actionBtn, { backgroundColor: '#fce4ec' }]}>
                <SymbolView name={{ ios: 'xmark', android: 'close', web: 'close' }} size={16} tintColor="#c62828" />
                <ThemedText style={{ color: '#c62828', fontSize: 13, fontWeight: '600' }}>Reject</ThemedText>
              </Pressable>
              <Pressable onPress={() => handleAction(item.id, 'approve')} style={[styles.actionBtn, { backgroundColor: '#e8f5e9' }]}>
                <SymbolView name={{ ios: 'checkmark', android: 'check', web: 'check' }} size={16} tintColor="#2e7d32" />
                <ThemedText style={{ color: '#2e7d32', fontSize: 13, fontWeight: '600' }}>Approve</ThemedText>
              </Pressable>
            </View>
          </ThemedView>
        )}
        ListEmptyComponent={
          <ThemedView style={styles.empty}>
            <SymbolView name={{ ios: 'checkmark.circle', android: 'check_circle', web: 'check_circle' }} size={40} tintColor={theme.textSecondary} />
            <ThemedText themeColor="textSecondary" style={{ textAlign: 'center', marginTop: Spacing.two }}>No pending verification requests</ThemedText>
          </ThemedView>
        }
      />

      <Modal visible={confirmModal.visible} animationType="fade" transparent>
        <Pressable style={styles.modalOverlay} onPress={() => setConfirmModal({ visible: false, id: '', action: 'approve' })}>
          <Pressable onPress={e => e.stopPropagation()}>
            <ThemedView type="backgroundElement" style={styles.confirmModal}>
              <SymbolView name={{ ios: confirmModal.action === 'approve' ? 'checkmark.circle' : 'xmark.circle' }} size={36} tintColor={confirmModal.action === 'approve' ? '#2e7d32' : '#c62828'} />
              <ThemedText style={{ fontWeight: '600', fontSize: 17, marginTop: Spacing.two }}>
                {confirmModal.action === 'approve' ? 'Approve Verification' : 'Reject Verification'}
              </ThemedText>
              <ThemedText themeColor="textSecondary" style={{ fontSize: 13, textAlign: 'center', marginTop: Spacing.one }}>
                {confirmModal.action === 'approve' ? 'This will verify the citizen identity.' : 'This will reject the verification request.'}
              </ThemedText>
              <View style={styles.confirmRow}>
                <Pressable onPress={() => setConfirmModal({ visible: false, id: '', action: 'approve' })} style={[styles.confirmBtn, { backgroundColor: theme.backgroundSelected }]}>
                  <ThemedText style={{ fontWeight: '600' }}>Cancel</ThemedText>
                </Pressable>
                <Pressable onPress={confirmAction} style={[styles.confirmBtn, { backgroundColor: confirmModal.action === 'approve' ? '#2e7d32' : '#c62828' }]}>
                  <ThemedText style={{ color: 'white', fontWeight: '600' }}>{confirmModal.action === 'approve' ? 'Approve' : 'Reject'}</ThemedText>
                </Pressable>
              </View>
            </ThemedView>
          </Pressable>
        </Pressable>
      </Modal>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1 },
  listContent: { paddingBottom: BottomTabInset + Spacing.three, maxWidth: MaxContentWidth, width: '100%', alignSelf: 'center' },
  header: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', paddingHorizontal: Spacing.four, paddingTop: Spacing.three, paddingBottom: Spacing.three },
  card: { marginHorizontal: Spacing.four, marginBottom: Spacing.two, padding: Spacing.three, borderRadius: Spacing.three, gap: Spacing.two },
  cardRow: { flexDirection: 'row', alignItems: 'center', gap: Spacing.three },
  avatar: { width: 44, height: 44, borderRadius: 22, justifyContent: 'center', alignItems: 'center' },
  cardInfo: { flex: 1, gap: 1 },
  actionRow: { flexDirection: 'row', gap: Spacing.two, marginTop: Spacing.one },
  actionBtn: { flex: 1, flexDirection: 'row', justifyContent: 'center', alignItems: 'center', paddingVertical: Spacing.two, borderRadius: Spacing.two, gap: Spacing.one },
  empty: { padding: Spacing.six, alignItems: 'center' },
  modalOverlay: { flex: 1, backgroundColor: 'rgba(0,0,0,0.4)', justifyContent: 'center', alignItems: 'center' },
  confirmModal: { padding: Spacing.four, borderRadius: Spacing.four, alignItems: 'center', minWidth: 280, maxWidth: 340 },
  confirmRow: { flexDirection: 'row', gap: Spacing.two, marginTop: Spacing.four, width: '100%' },
  confirmBtn: { flex: 1, paddingVertical: Spacing.three, borderRadius: Spacing.two, alignItems: 'center' },
});
