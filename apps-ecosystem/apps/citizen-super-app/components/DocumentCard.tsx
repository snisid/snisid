import { View, Text, StyleSheet, Pressable } from 'react-native';
import { useColorScheme } from './useColorScheme';

import type { Document } from '@/lib/api';

interface DocumentCardProps {
  document: Document;
  onPress?: (doc: Document) => void;
}

const TYPE_ICONS: Record<string, string> = {
  birth_certificate: '📜',
  national_id: '🪪',
  passport: '🛂',
  drivers_license: '🚗',
};

const STATUS_STYLES: Record<string, { bg: string; text: string; label: string }> = {
  valid: { bg: '#DCFCE7', text: '#166534', label: 'Valid' },
  expired: { bg: '#FEF3C7', text: '#92400E', label: 'Expired' },
  pending: { bg: '#DBEAFE', text: '#1E40AF', label: 'Pending' },
};

export default function DocumentCard({ document: doc, onPress }: DocumentCardProps) {
  const colorScheme = useColorScheme();
  const isDark = colorScheme === 'dark';
  const statusInfo = STATUS_STYLES[doc.status] ?? STATUS_STYLES.pending;

  return (
    <Pressable
      onPress={() => onPress?.(doc)}
      style={({ pressed }) => [
        styles.card,
        isDark && styles.cardDark,
        pressed && styles.pressed,
      ]}>
      <View style={styles.iconContainer}>
        <Text style={styles.icon}>{TYPE_ICONS[doc.type] ?? '📄'}</Text>
      </View>

      <View style={styles.content}>
        <Text style={[styles.name, isDark && styles.textDark]}>{doc.name}</Text>
        <Text style={[styles.number, isDark && styles.textDark]}>
          {doc.documentNumber}
        </Text>
        <Text style={[styles.date, isDark && styles.textDark]}>
          Issued: {new Date(doc.issueDate).toLocaleDateString()}
          {' · '}Exp: {new Date(doc.expirationDate).toLocaleDateString()}
        </Text>
        <Text style={[styles.issuer, isDark && styles.textDark]}>{doc.issuer}</Text>
      </View>

      <View style={[styles.statusBadge, { backgroundColor: statusInfo.bg }]}>
        <Text style={[styles.statusText, { color: statusInfo.text }]}>{statusInfo.label}</Text>
      </View>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  card: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 16,
    marginHorizontal: 16,
    marginVertical: 4,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.05,
    shadowRadius: 4,
    elevation: 2,
  },
  cardDark: {
    backgroundColor: '#1C1C1E',
  },
  pressed: {
    opacity: 0.7,
  },
  iconContainer: {
    width: 48,
    height: 48,
    borderRadius: 12,
    backgroundColor: '#F0F0F5',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 12,
  },
  icon: {
    fontSize: 24,
  },
  content: {
    flex: 1,
    gap: 2,
  },
  name: {
    fontSize: 16,
    fontWeight: '600',
    color: '#000',
  },
  number: {
    fontSize: 12,
    color: '#555',
    fontFamily: 'monospace',
  },
  date: {
    fontSize: 11,
    color: '#888',
    marginTop: 2,
  },
  issuer: {
    fontSize: 11,
    color: '#888',
  },
  statusBadge: {
    paddingHorizontal: 8,
    paddingVertical: 3,
    borderRadius: 8,
    marginLeft: 8,
  },
  statusText: {
    fontSize: 11,
    fontWeight: '700',
  },
  textDark: {
    color: '#fff',
  },
});
