import { useState, useCallback } from 'react';
import { ScrollView, Pressable, RefreshControl, StyleSheet, View, Modal } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { SymbolView } from 'expo-symbols';

import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';

import { MaxContentWidth, Spacing } from '@/constants/theme';
import { useTheme } from '@/hooks/use-theme';

type TimelineEvent = { id: string; event: string; timestamp: string; severity: 'critical' | 'warning' | 'info' };

const TIMELINE: TimelineEvent[] = [
  { id: '1', event: 'Identity verified successfully', timestamp: '2 days ago', severity: 'info' },
  { id: '2', event: 'Biometric enrollment completed', timestamp: '2 days ago', severity: 'info' },
  { id: '3', event: 'Document #DOC-4452 uploaded', timestamp: '3 days ago', severity: 'info' },
  { id: '4', event: 'Registration submitted', timestamp: '5 days ago', severity: 'info' },
];

const SEVERITY_CONFIG = { critical: { color: '#c62828', bg: '#fce4ec' }, warning: { color: '#e65100', bg: '#fff3e0' }, info: { color: '#1565c0', bg: '#e3f2fd' } };

export default function CitizenDetailScreen() {
  const theme = useTheme();
  const router = useRouter();
  const { id } = useLocalSearchParams<{ id: string }>();
  const [refreshing, setRefreshing] = useState(false);
  const [actionModal, setActionModal] = useState<{ visible: boolean; action: string }>({ visible: false, action: '' });

  const onRefresh = useCallback(() => {
    setRefreshing(true);
    setTimeout(() => setRefreshing(false), 1500);
  }, []);

  const handleAction = (action: string) => setActionModal({ visible: true, action });
  const confirmAction = () => setActionModal({ visible: false, action: '' });

  const InfoRow = ({ label, value }: { label: string; value: string }) => (
    <View style={styles.infoRow}>
      <ThemedText themeColor="textSecondary" style={{ fontSize: 12 }}>{label}</ThemedText>
      <ThemedText style={{ fontSize: 14 }}>{value}</ThemedText>
    </View>
  );

  return (
    <SafeAreaView style={[styles.safeArea, { backgroundColor: theme.background }]}>
      <ScrollView
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
        contentContainerStyle={styles.scrollContent}
      >
        <View style={styles.headerRow}>
          <Pressable onPress={() => router.back()} style={styles.backBtn}>
            <SymbolView name={{ ios: 'chevron.left', android: 'arrow_back', web: 'arrow_back' }} size={20} tintColor={theme.text} />
          </Pressable>
          <ThemedText type="subtitle" style={{ fontSize: 22, lineHeight: 28 }}>Citizen Profile</ThemedText>
        </View>

        <ThemedView type="backgroundElement" style={styles.section}>
          <View style={styles.avatarRow}>
            <View style={[styles.avatar, { backgroundColor: theme.backgroundSelected }]}>
              <ThemedText style={{ fontWeight: '700', fontSize: 24 }}>AS</ThemedText>
            </View>
            <View style={styles.avatarInfo}>
              <ThemedText style={{ fontWeight: '600', fontSize: 18 }}>Abubakar Sani</ThemedText>
              <View style={[styles.statusBadge, { backgroundColor: '#e8f5e9' }]}>
                <ThemedText style={{ color: '#2e7d32', fontSize: 11, fontWeight: '700', textTransform: 'uppercase' }}>Active</ThemedText>
              </View>
            </View>
          </View>
          <InfoRow label="National ID (NNU)" value="SN-88472" />
          <InfoRow label="Date of Birth" value="1987-04-15" />
          <InfoRow label="Gender" value="Male" />
          <InfoRow label="Phone" value="+234 802 123 4567" />
          <InfoRow label="Email" value="a.sani@example.com" />
          <InfoRow label="State of Origin" value="Kaduna" />
          <InfoRow label="LGA" value="Zaria" />
          <InfoRow label="Agency" value="NIA HQ" />
        </ThemedView>

        <ThemedView type="backgroundElement" style={styles.section}>
          <View style={styles.sectionHeader}>
            <SymbolView name={{ ios: 'doc.text', android: 'description', web: 'description' }} size={18} tintColor={theme.text} />
            <ThemedText style={{ fontWeight: '600', fontSize: 15 }}>Documents (3)</ThemedText>
          </View>
          {['Birth Certificate (DOC-4452)', 'Passport Photo (DOC-4453)', 'Proof of Address (DOC-4454)'].map((doc, i) => (
            <View key={i} style={styles.docRow}>
              <SymbolView name={{ ios: 'doc', android: 'description', web: 'description' }} size={16} tintColor={theme.textSecondary} />
              <ThemedText style={{ fontSize: 13, flex: 1 }}>{doc}</ThemedText>
              <ThemedText themeColor="textSecondary" style={{ fontSize: 11 }}>Verified</ThemedText>
            </View>
          ))}
        </ThemedView>

        <ThemedView type="backgroundElement" style={styles.section}>
          <View style={styles.sectionHeader}>
            <SymbolView name={{ ios: 'touchid', android: 'fingerprint', web: 'fingerprint' }} size={18} tintColor={theme.text} />
            <ThemedText style={{ fontWeight: '600', fontSize: 15 }}>Biometric Enrollment</ThemedText>
          </View>
          <InfoRow label="Fingerprints" value="10/10 enrolled" />
          <InfoRow label="Facial Recognition" value="Enrolled" />
          <InfoRow label="Iris Scan" value="Not enrolled" />
          <InfoRow label="Signature" value="Captured" />
        </ThemedView>

        <ThemedView type="backgroundElement" style={styles.section}>
          <View style={styles.sectionHeader}>
            <SymbolView name={{ ios: 'clock.arrow.circlepath', android: 'history', web: 'history' }} size={18} tintColor={theme.text} />
            <ThemedText style={{ fontWeight: '600', fontSize: 15 }}>Verification History</ThemedText>
          </View>
          {TIMELINE.map(evt => {
            const cfg = SEVERITY_CONFIG[evt.severity];
            return (
              <View key={evt.id} style={styles.timelineRow}>
                <View style={[styles.timelineDot, { backgroundColor: cfg.color }]} />
                <View style={styles.timelineContent}>
                  <ThemedText style={{ fontSize: 13 }}>{evt.event}</ThemedText>
                  <ThemedText themeColor="textSecondary" style={{ fontSize: 11 }}>{evt.timestamp}</ThemedText>
                </View>
              </View>
            );
          })}
        </ThemedView>

        <View style={styles.actionRow}>
          <Pressable onPress={() => handleAction('verify')} style={[styles.actionBtn, { backgroundColor: '#e8f5e9' }]}>
            <SymbolView name={{ ios: 'checkmark.circle', android: 'check_circle', web: 'check_circle' }} size={18} tintColor="#2e7d32" />
            <ThemedText style={{ color: '#2e7d32', fontWeight: '600', fontSize: 13 }}>Verify</ThemedText>
          </Pressable>
          <Pressable onPress={() => handleAction('suspend')} style={[styles.actionBtn, { backgroundColor: '#fff3e0' }]}>
            <SymbolView name={{ ios: 'pause.circle', android: 'pause', web: 'pause' }} size={18} tintColor="#e65100" />
            <ThemedText style={{ color: '#e65100', fontWeight: '600', fontSize: 13 }}>Suspend</ThemedText>
          </Pressable>
          <Pressable onPress={() => handleAction('revoke')} style={[styles.actionBtn, { backgroundColor: '#fce4ec' }]}>
            <SymbolView name={{ ios: 'xmark.circle', android: 'cancel', web: 'cancel' }} size={18} tintColor="#c62828" />
            <ThemedText style={{ color: '#c62828', fontWeight: '600', fontSize: 13 }}>Revoke</ThemedText>
          </Pressable>
        </View>

        <Modal visible={actionModal.visible} animationType="fade" transparent>
          <Pressable style={styles.modalOverlay} onPress={() => setActionModal({ visible: false, action: '' })}>
            <Pressable onPress={e => e.stopPropagation()}>
              <ThemedView type="backgroundElement" style={styles.confirmModal}>
                <ThemedText style={{ fontWeight: '600', fontSize: 17, textTransform: 'capitalize' }}>{actionModal.action} Identity</ThemedText>
                <ThemedText themeColor="textSecondary" style={{ fontSize: 13, textAlign: 'center', marginTop: Spacing.one }}>
                  Are you sure you want to {actionModal.action} citizen SN-88472?
                </ThemedText>
                <View style={styles.confirmRow}>
                  <Pressable onPress={() => setActionModal({ visible: false, action: '' })} style={[styles.confirmBtn, { backgroundColor: theme.backgroundSelected }]}>
                    <ThemedText style={{ fontWeight: '600' }}>Cancel</ThemedText>
                  </Pressable>
                  <Pressable onPress={confirmAction} style={[styles.confirmBtn, { backgroundColor: actionModal.action === 'verify' ? '#2e7d32' : actionModal.action === 'suspend' ? '#e65100' : '#c62828' }]}>
                    <ThemedText style={{ color: 'white', fontWeight: '600', textTransform: 'capitalize' }}>{actionModal.action}</ThemedText>
                  </Pressable>
                </View>
              </ThemedView>
            </Pressable>
          </Pressable>
        </Modal>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1 },
  scrollContent: { paddingBottom: Spacing.six, maxWidth: MaxContentWidth, width: '100%', alignSelf: 'center', paddingHorizontal: Spacing.four },
  headerRow: { flexDirection: 'row', alignItems: 'center', gap: Spacing.two, paddingTop: Spacing.three, paddingBottom: Spacing.three },
  backBtn: { padding: Spacing.one },
  section: { padding: Spacing.three, borderRadius: Spacing.three, marginBottom: Spacing.three, gap: Spacing.two },
  sectionHeader: { flexDirection: 'row', alignItems: 'center', gap: Spacing.two },
  avatarRow: { flexDirection: 'row', alignItems: 'center', gap: Spacing.three, marginBottom: Spacing.two },
  avatar: { width: 56, height: 56, borderRadius: 28, justifyContent: 'center', alignItems: 'center' },
  avatarInfo: { gap: Spacing.one },
  statusBadge: { paddingHorizontal: Spacing.two, paddingVertical: 2, borderRadius: Spacing.one, alignSelf: 'flex-start' },
  infoRow: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', paddingVertical: Spacing.one },
  docRow: { flexDirection: 'row', alignItems: 'center', gap: Spacing.two, paddingVertical: Spacing.one },
  timelineRow: { flexDirection: 'row', gap: Spacing.two, paddingVertical: Spacing.one },
  timelineDot: { width: 8, height: 8, borderRadius: 4, marginTop: 6 },
  timelineContent: { flex: 1, gap: 1 },
  actionRow: { flexDirection: 'row', gap: Spacing.two, marginVertical: Spacing.four },
  actionBtn: { flex: 1, flexDirection: 'row', justifyContent: 'center', alignItems: 'center', paddingVertical: Spacing.three, borderRadius: Spacing.two, gap: Spacing.one },
  modalOverlay: { flex: 1, backgroundColor: 'rgba(0,0,0,0.4)', justifyContent: 'center', alignItems: 'center' },
  confirmModal: { padding: Spacing.four, borderRadius: Spacing.four, alignItems: 'center', minWidth: 280, maxWidth: 340 },
  confirmRow: { flexDirection: 'row', gap: Spacing.two, marginTop: Spacing.four, width: '100%' },
  confirmBtn: { flex: 1, paddingVertical: Spacing.three, borderRadius: Spacing.two, alignItems: 'center' },
});
