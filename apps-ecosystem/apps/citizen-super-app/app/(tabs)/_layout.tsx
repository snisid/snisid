import { Tabs } from 'expo-router';
import { Platform, Pressable } from 'react-native';
import { SymbolView } from 'expo-symbols';

import Colors from '@/constants/Colors';
import { useColorScheme } from '@/components/useColorScheme';
import { useClientOnlyValue } from '@/components/useClientOnlyValue';
import { useAuth } from '@/lib/auth';

export default function TabLayout() {
  const colorScheme = useColorScheme();
  const { isAuthenticated } = useAuth();

  return (
    <Tabs
      screenOptions={{
        tabBarActiveTintColor: Colors[colorScheme].tint,
        tabBarInactiveTintColor: Colors[colorScheme].tabIconDefault,
        headerShown: useClientOnlyValue(false, true),
        tabBarStyle: {
          borderTopColor: Colors[colorScheme].tabIconDefault + '40',
        },
      }}>
      <Tabs.Screen
        name="index"
        options={{
          title: 'Home',
          tabBarIcon: ({ color, size }) => (
            <SymbolView
              name={{ ios: 'house.fill', android: 'home', web: 'home' }}
              tintColor={color}
              size={size}
            />
          ),
          href: isAuthenticated ? '/(tabs)' : null,
        }}
      />
      <Tabs.Screen
        name="two"
        options={{
          title: 'Documents',
          tabBarIcon: ({ color, size }) => (
            <SymbolView
              name={{ ios: 'doc.text.fill', android: 'description', web: 'description' }}
              tintColor={color}
              size={size}
            />
          ),
          href: isAuthenticated ? '/(tabs)/two' : null,
        }}
      />
      <Tabs.Screen
        name="notifications"
        options={{
          title: 'Activity',
          tabBarIcon: ({ color, size }) => (
            <SymbolView
              name={{ ios: 'bell.fill', android: 'notifications', web: 'notifications' }}
              tintColor={color}
              size={size}
            />
          ),
          href: isAuthenticated ? '/(tabs)/notifications' : null,
        }}
      />
      <Tabs.Screen
        name="support"
        options={{
          title: 'Support',
          tabBarIcon: ({ color, size }) => (
            <SymbolView
              name={{ ios: 'questionmark.circle.fill', android: 'support', web: 'support' }}
              tintColor={color}
              size={size}
            />
          ),
          href: isAuthenticated ? '/(tabs)/support' : null,
        }}
      />
      <Tabs.Screen
        name="settings"
        options={{
          title: 'Settings',
          tabBarIcon: ({ color, size }) => (
            <SymbolView
              name={{ ios: 'gearshape.fill', android: 'settings', web: 'settings' }}
              tintColor={color}
              size={size}
            />
          ),
          href: isAuthenticated ? '/(tabs)/settings' : null,
        }}
      />
    </Tabs>
  );
}
