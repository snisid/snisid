import { View, Text, StyleSheet } from 'react-native';
import { useColorScheme } from './useColorScheme';

export type ActivityStatus = 'success' | 'warning' | 'error' | 'info';

interface ActivityCardProps {
  icon: string;
  title: string;
  description: string;
  timestamp: string;
  status?: ActivityStatus;
}

const STATUS_INDICATORS: Record<ActivityStatus, { color: string; bg: string }> = {
  success: { color: '#166534', bg: '#DCFCE7' },
  warning: { color: '#92400E', bg: '#FEF3C7' },
  error: { color: '#991B1B', bg: '#FEE2E2' },
  info: { color: '#1E40AF', bg: '#DBEAFE' },
};

function formatRelativeTime(timestamp: string): string {
  const now = Date.now();
  const date = new Date(timestamp).getTime();
  const diffMs = now - date;
  const diffMin = Math.floor(diffMs / 60000);
  if (diffMin < 1) return 'Just now';
  if (diffMin < 60) return `${diffMin}m ago`;
  const diffHrs = Math.floor(diffMin / 60);
  if (diffHrs < 24) return `${diffHrs}h ago`;
  const diffDays = Math.floor(diffHrs / 24);
  if (diffDays < 7) return `${diffDays}d ago`;
  return new Date(timestamp).toLocaleDateString();
}

export default function ActivityCard({
  icon,
  title,
  description,
  timestamp,
  status = 'info',
}: ActivityCardProps) {
  const colorScheme = useColorScheme();
  const isDark = colorScheme === 'dark';
  const indicator = STATUS_INDICATORS[status];

  return (
    <View style={[styles.card, isDark && styles.cardDark]}>
      <View style={[styles.iconContainer, { backgroundColor: indicator.bg }]}>
        <Text style={styles.icon}>{icon}</Text>
      </View>
      <View style={styles.content}>
        <Text style={[styles.title, isDark && styles.textDark]} numberOfLines={1}>
          {title}
        </Text>
        <Text style={[styles.description, isDark && styles.textDark]} numberOfLines={2}>
          {description}
        </Text>
        <Text style={styles.timestamp}>{formatRelativeTime(timestamp)}</Text>
      </View>
      <View style={[styles.statusDot, { backgroundColor: indicator.color }]} />
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#fff',
    borderRadius: 14,
    padding: 14,
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
  iconContainer: {
    width: 44,
    height: 44,
    borderRadius: 12,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 12,
  },
  icon: {
    fontSize: 22,
  },
  content: {
    flex: 1,
    gap: 2,
  },
  title: {
    fontSize: 15,
    fontWeight: '600',
    color: '#000',
  },
  description: {
    fontSize: 13,
    color: '#666',
    lineHeight: 18,
    marginTop: 1,
  },
  timestamp: {
    fontSize: 11,
    color: '#999',
    marginTop: 4,
  },
  statusDot: {
    width: 8,
    height: 8,
    borderRadius: 4,
    marginLeft: 8,
  },
  textDark: {
    color: '#fff',
  },
});
