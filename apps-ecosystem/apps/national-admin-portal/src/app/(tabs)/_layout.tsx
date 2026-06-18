import { Tabs } from 'expo-router';
import { SymbolView } from 'expo-symbols';
import { useColorScheme } from 'react-native';

import { Colors } from '@/constants/theme';

export default function TabLayout() {
  const scheme = useColorScheme();
  const colors = Colors[scheme === 'unspecified' ? 'light' : scheme];

  return (
    <Tabs
      screenOptions={{
        tabBarActiveTintColor: colors.text,
        tabBarInactiveTintColor: colors.textSecondary,
        tabBarStyle: { backgroundColor: colors.background, borderTopColor: colors.backgroundElement, borderTopWidth: 1 },
        headerStyle: { backgroundColor: colors.background },
        headerTintColor: colors.text,
      }}>
      <Tabs.Screen
        name="index"
        options={{
          title: 'Dashboard',
          tabBarIcon: ({ color, size }) => (
            <SymbolView name={{ ios: 'square.grid.2x2', android: 'grid_view', web: 'grid_view' }} size={size} tintColor={color} />
          ),
        }}
      />
      <Tabs.Screen
        name="citizens"
        options={{
          title: 'Citizens',
          tabBarIcon: ({ color, size }) => (
            <SymbolView name={{ ios: 'person.2', android: 'people', web: 'people' }} size={size} tintColor={color} />
          ),
        }}
      />
      <Tabs.Screen
        name="verification"
        options={{
          title: 'Verification',
          tabBarIcon: ({ color, size }) => (
            <SymbolView name={{ ios: 'checkmark.shield', android: 'verified_user', web: 'verified_user' }} size={size} tintColor={color} />
          ),
        }}
      />
      <Tabs.Screen
        name="security"
        options={{
          title: 'Security',
          tabBarIcon: ({ color, size }) => (
            <SymbolView name={{ ios: 'lock.shield', android: 'security', web: 'security' }} size={size} tintColor={color} />
          ),
        }}
      />
      <Tabs.Screen
        name="settings"
        options={{
          title: 'Settings',
          tabBarIcon: ({ color, size }) => (
            <SymbolView name={{ ios: 'gearshape', android: 'settings', web: 'settings' }} size={size} tintColor={color} />
          ),
        }}
      />
    </Tabs>
  );
}
