import { useState } from 'react';
import { Alert, Pressable, StyleSheet, View } from 'react-native';
import { SymbolView } from 'expo-symbols';

import { ThemedText } from './themed-text';
import { ThemedView } from './themed-view';

import { Spacing } from '@/constants/theme';
import { useTheme } from '@/hooks/use-theme';

type VerificationRequestProps = {
  citizenName: string;
  nnu: string;
  requestType: string;
  submittedAt: string;
  documentType: string;
  onAccept: () => void;
  onReject: () => void;
};

export function VerificationRequest({
  citizenName, nnu, requestType, submittedAt, documentType, onAccept, onReject,
}: VerificationRequestProps) {
  const [expanded, setExpanded] = useState(false);
  const theme = useTheme();

  const handleAccept = () => {
    Alert.alert('Confirm Verification', `Verify identity for ${citizenName} (${nnu})?`, [
      { text: 'Cancel', style: 'cancel' },
      { text: 'Verify', onPress: onAccept },
    ]);
  };

  const handleReject = () => {
    Alert.alert('Reject Request', `Reject verification for ${citizenName} (${nnu})?`, [
      { text: 'Cancel', style: 'cancel' },
      { text: 'Reject', style: 'destructive', onPress: onReject },
    ]);
  };

  return (
    <ThemedView type="backgroundElement" style={styles.card}>
      <Pressable onPress={() => setExpanded(!expanded)} style={styles.header}>
        <View style={styles.headerInfo}>
          <ThemedText style={{ fontWeight: '600' }}>{citizenName}</ThemedText>
          <ThemedText themeColor="textSecondary" style={{ fontSize: 13 }}>NNU: {nnu}</ThemedText>
        </View>
        <SymbolView
          name={{ ios: 'chevron.right', android: 'chevron_right', web: 'chevron_right' }}
          size={14}
          tintColor={theme.textSecondary}
        />
      </Pressable>

      {expanded && (
        <View style={styles.details}>
          <View style={styles.detailRow}>
            <ThemedText themeColor="textSecondary" style={{ fontSize: 13 }}>Request</ThemedText>
            <ThemedText style={{ fontSize: 13 }}>{requestType}</ThemedText>
          </View>
          <View style={styles.detailRow}>
            <ThemedText themeColor="textSecondary" style={{ fontSize: 13 }}>Document</ThemedText>
            <ThemedText style={{ fontSize: 13 }}>{documentType}</ThemedText>
          </View>
          <View style={styles.detailRow}>
            <ThemedText themeColor="textSecondary" style={{ fontSize: 13 }}>Submitted</ThemedText>
            <ThemedText style={{ fontSize: 13 }}>{submittedAt}</ThemedText>
          </View>
          <View style={styles.actions}>
            <Pressable onPress={handleReject} style={styles.rejectBtn}>
              <ThemedText style={{ color: '#c62828', fontWeight: '600', fontSize: 13 }}>Reject</ThemedText>
            </Pressable>
            <Pressable onPress={handleAccept} style={styles.acceptBtn}>
              <ThemedText style={{ color: '#2e7d32', fontWeight: '600', fontSize: 13 }}>Verify</ThemedText>
            </Pressable>
          </View>
        </View>
      )}
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  card: {
    borderRadius: Spacing.three,
    overflow: 'hidden',
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    padding: Spacing.three,
  },
  headerInfo: {
    flex: 1,
    gap: 2,
  },
  details: {
    paddingHorizontal: Spacing.three,
    paddingBottom: Spacing.three,
    gap: Spacing.two,
  },
  detailRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  actions: {
    flexDirection: 'row',
    gap: Spacing.two,
    marginTop: Spacing.one,
  },
  rejectBtn: {
    flex: 1,
    paddingVertical: Spacing.two,
    borderRadius: Spacing.two,
    alignItems: 'center',
    backgroundColor: '#fce4ec',
  },
  acceptBtn: {
    flex: 1,
    paddingVertical: Spacing.two,
    borderRadius: Spacing.two,
    alignItems: 'center',
    backgroundColor: '#e8f5e9',
  },
});
