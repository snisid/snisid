import { View, Text, StyleSheet, Switch } from 'react-native';
import { useColorScheme } from './useColorScheme';

interface Permission {
  id: string;
  name: string;
  icon: string;
  description: string;
  granted: boolean;
}

interface PermissionsListProps {
  permissions: Permission[];
  onToggle: (id: string, newValue: boolean) => void;
  title?: string;
}

export default function PermissionsList({
  permissions,
  onToggle,
  title,
}: PermissionsListProps) {
  const colorScheme = useColorScheme();
  const isDark = colorScheme === 'dark';

  if (permissions.length === 0) {
    return (
      <View style={styles.emptyContainer}>
        <Text style={styles.emptyText}>No permissions granted</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      {title && <Text style={styles.sectionTitle}>{title}</Text>}
      {permissions.map((perm, index) => (
        <View
          key={perm.id}
          style={[
            styles.row,
            isDark && styles.rowDark,
            index === 0 && styles.rowFirst,
            index === permissions.length - 1 && styles.rowLast,
          ]}>
          <View style={styles.iconContainer}>
            <Text style={styles.icon}>{perm.icon}</Text>
          </View>
          <View style={styles.info}>
            <Text style={[styles.name, isDark && styles.textDark]}>{perm.name}</Text>
            <Text style={[styles.description, isDark && styles.textDark]}>
              {perm.description}
            </Text>
          </View>
          <Switch
            value={perm.granted}
            onValueChange={(value) => onToggle(perm.id, value)}
            trackColor={{ false: '#E5E5EA', true: '#0033a080' }}
            thumbColor={perm.granted ? '#0033a0' : '#f4f3f4'}
          />
        </View>
      ))}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    marginHorizontal: 16,
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
  row: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#fff',
    padding: 14,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: '#E5E5EA',
  },
  rowDark: {
    backgroundColor: '#1C1C1E',
    borderBottomColor: '#38383A',
  },
  rowFirst: {
    borderTopLeftRadius: 14,
    borderTopRightRadius: 14,
  },
  rowLast: {
    borderBottomLeftRadius: 14,
    borderBottomRightRadius: 14,
    borderBottomWidth: 0,
  },
  iconContainer: {
    width: 36,
    height: 36,
    borderRadius: 10,
    backgroundColor: '#F0F0F5',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 12,
  },
  icon: {
    fontSize: 18,
  },
  info: {
    flex: 1,
    marginRight: 8,
  },
  name: {
    fontSize: 15,
    fontWeight: '600',
    color: '#000',
  },
  description: {
    fontSize: 12,
    color: '#888',
    marginTop: 1,
  },
  emptyContainer: {
    padding: 24,
    alignItems: 'center',
  },
  emptyText: {
    fontSize: 14,
    color: '#999',
  },
  textDark: {
    color: '#fff',
  },
});
