import { View, Text, StyleSheet, Image } from 'react-native';
import { useColorScheme } from './useColorScheme';

import Colors from '@/constants/Colors';
import type { Identity } from '@/lib/api';

interface IdentityCardProps {
  identity: Identity;
}

const STATUS_COLORS: Record<string, { bg: string; text: string; label: string }> = {
  active: { bg: '#DCFCE7', text: '#166534', label: 'Active' },
  expired: { bg: '#FEF3C7', text: '#92400E', label: 'Expired' },
  suspended: { bg: '#FEE2E2', text: '#991B1B', label: 'Suspended' },
  pending: { bg: '#DBEAFE', text: '#1E40AF', label: 'Pending' },
};

export default function IdentityCard({ identity }: IdentityCardProps) {
  const colorScheme = useColorScheme();
  const isDark = colorScheme === 'dark';
  const statusInfo = STATUS_COLORS[identity.status] ?? STATUS_COLORS.pending;

  return (
    <View style={[styles.card, isDark && styles.cardDark]}>
      <View style={styles.header}>
        <View style={[styles.statusBadge, { backgroundColor: statusInfo.bg }]}>
          <Text style={[styles.statusText, { color: statusInfo.text }]}>{statusInfo.label}</Text>
        </View>
      </View>

      <View style={styles.body}>
        <View style={styles.photoContainer}>
          {identity.photoUrl ? (
            <Image source={{ uri: identity.photoUrl }} style={styles.photo} />
          ) : (
            <View style={[styles.photoPlaceholder, isDark && styles.photoPlaceholderDark]}>
              <Text style={styles.photoInitials}>
                {identity.fullName
                  .split(' ')
                  .map((n: string) => n[0])
                  .slice(0, 2)
                  .join('')
                  .toUpperCase()}
              </Text>
            </View>
          )}
        </View>

        <View style={styles.info}>
          <Text style={[styles.nnu, isDark && styles.textDark]}>NNU: {identity.nnu}</Text>
          <Text style={[styles.name, isDark && styles.textDark]}>{identity.fullName}</Text>
          <Text style={[styles.detail, isDark && styles.textDark]}>
            DOB: {new Date(identity.dateOfBirth).toLocaleDateString()}
          </Text>
          <Text style={[styles.detail, isDark && styles.textDark]}>
            Gender: {identity.gender}
          </Text>
          <Text style={[styles.detail, isDark && styles.textDark]}>
            Nationality: {identity.nationality}
          </Text>
        </View>
      </View>

      <View style={[styles.footer, isDark && styles.footerDark]}>
        <Text style={[styles.detail, isDark && styles.textDark]}>
          Expires: {new Date(identity.expirationDate).toLocaleDateString()}
        </Text>
      </View>

      <View style={styles.qrPlaceholder}>
        <View style={[styles.qrBox, isDark && styles.qrBoxDark]}>
          <Text style={[styles.qrIcon, isDark && styles.textDark]}>◆◆</Text>
          <Text style={[styles.qrLabel, isDark && styles.textDark]}>QR Code</Text>
        </View>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: '#fff',
    borderRadius: 16,
    padding: 20,
    marginHorizontal: 16,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 8,
    elevation: 4,
  },
  cardDark: {
    backgroundColor: '#1C1C1E',
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
    marginBottom: 16,
  },
  statusBadge: {
    paddingHorizontal: 12,
    paddingVertical: 4,
    borderRadius: 12,
  },
  statusText: {
    fontSize: 12,
    fontWeight: '700',
    textTransform: 'uppercase',
    letterSpacing: 0.5,
  },
  body: {
    flexDirection: 'row',
    gap: 16,
    marginBottom: 16,
  },
  photoContainer: {
    width: 80,
    height: 100,
    borderRadius: 8,
    overflow: 'hidden',
  },
  photo: {
    width: 80,
    height: 100,
  },
  photoPlaceholder: {
    width: 80,
    height: 100,
    backgroundColor: '#E0E0E0',
    justifyContent: 'center',
    alignItems: 'center',
    borderRadius: 8,
  },
  photoPlaceholderDark: {
    backgroundColor: '#2C2C2E',
  },
  photoInitials: {
    fontSize: 28,
    fontWeight: '700',
    color: '#666',
  },
  info: {
    flex: 1,
    gap: 4,
  },
  nnu: {
    fontSize: 14,
    fontWeight: '700',
    color: '#0033a0',
    fontFamily: 'monospace',
  },
  name: {
    fontSize: 18,
    fontWeight: '600',
    color: '#000',
  },
  detail: {
    fontSize: 13,
    color: '#555',
  },
  footer: {
    borderTopWidth: 1,
    borderTopColor: '#E5E5EA',
    paddingTop: 12,
    alignItems: 'center',
  },
  footerDark: {
    borderTopColor: '#38383A',
  },
  qrPlaceholder: {
    alignItems: 'center',
    marginTop: 16,
  },
  qrBox: {
    width: 100,
    height: 100,
    backgroundColor: '#F5F5F5',
    borderRadius: 12,
    justifyContent: 'center',
    alignItems: 'center',
    borderWidth: 2,
    borderColor: '#E0E0E0',
    borderStyle: 'dashed',
  },
  qrBoxDark: {
    backgroundColor: '#2C2C2E',
    borderColor: '#38383A',
  },
  qrIcon: {
    fontSize: 24,
    color: '#0033a0',
  },
  qrLabel: {
    fontSize: 10,
    color: '#666',
    marginTop: 4,
  },
  textDark: {
    color: '#fff',
  },
});
