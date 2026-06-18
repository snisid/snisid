import { Stack } from 'expo-router';
import { DarkTheme, DefaultTheme, ThemeProvider } from 'expo-router';
import { useColorScheme } from 'react-native';

import { AnimatedSplashOverlay } from '@/components/animated-icon';

export default function RootLayout() {
  const colorScheme = useColorScheme();
  return (
    <ThemeProvider value={colorScheme === 'dark' ? DarkTheme : DefaultTheme}>
      <AnimatedSplashOverlay />
      <Stack>
        <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
        <Stack.Screen name="citizen/[id]" options={{ presentation: 'modal', title: 'Citizen Details' }} />
        <Stack.Screen name="audit-log" options={{ presentation: 'modal', title: 'Audit Log' }} />
        <Stack.Screen name="reports" options={{ presentation: 'modal', title: 'Reports & Analytics' }} />
        <Stack.Screen name="alerts" options={{ presentation: 'modal', title: 'Alert Center' }} />
      </Stack>
    </ThemeProvider>
  );
}
