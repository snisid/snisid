import { useRef } from 'react';
import { TextInput, Pressable, StyleSheet, View } from 'react-native';
import { SymbolView } from 'expo-symbols';

import { ThemedView } from './themed-view';

import { Spacing } from '@/constants/theme';
import { useTheme } from '@/hooks/use-theme';

type SearchBarProps = {
  placeholder?: string;
  value: string;
  onChangeText: (text: string) => void;
  onFilterPress?: () => void;
};

export function SearchBar({ placeholder = 'Search...', value, onChangeText, onFilterPress }: SearchBarProps) {
  const theme = useTheme();
  const inputRef = useRef<TextInput>(null);

  return (
    <ThemedView type="backgroundElement" style={styles.container}>
      <SymbolView name={{ ios: 'magnifyingglass', android: 'search', web: 'search' }} size={18} tintColor={theme.textSecondary} />
      <TextInput
        ref={inputRef}
        style={[styles.input, { color: theme.text }]}
        placeholder={placeholder}
        placeholderTextColor={theme.textSecondary}
        value={value}
        onChangeText={onChangeText}
        autoCapitalize="none"
        autoCorrect={false}
      />
      {value.length > 0 && (
        <Pressable onPress={() => { onChangeText(''); inputRef.current?.focus(); }} style={styles.clearBtn} hitSlop={8}>
          <SymbolView name={{ ios: 'xmark.circle.fill', android: 'cancel', web: 'cancel' }} size={18} tintColor={theme.textSecondary} />
        </Pressable>
      )}
      {onFilterPress && (
        <Pressable onPress={onFilterPress} style={styles.filterBtn} hitSlop={8}>
          <SymbolView name={{ ios: 'line.3.horizontal.decrease', android: 'filter_list', web: 'filter_list' }} size={18} tintColor={theme.text} />
        </Pressable>
      )}
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: Spacing.three,
    paddingVertical: 10,
    borderRadius: Spacing.three,
    gap: Spacing.two,
  },
  input: {
    flex: 1,
    fontSize: 15,
    paddingVertical: 0,
  },
  clearBtn: {
    padding: Spacing.half,
  },
  filterBtn: {
    padding: Spacing.half,
    marginLeft: Spacing.one,
  },
});
