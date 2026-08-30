import { useState, useCallback } from 'react';
import { FlatList, Pressable, RefreshControl, StyleSheet, View, Modal } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter } from 'expo-router';

import { CitizenCard, type CitizenStatus } from '@/components/CitizenCard';
import { SearchBar } from '@/components/SearchBar';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';

import { BottomTabInset, MaxContentWidth, Spacing } from '@/constants/theme';
import { useTheme } from '@/hooks/use-theme';

type Citizen = {
  id: string;
  name: string;
  nnu: string;
  status: CitizenStatus;
  agency: string;
};

const MOCK_CITIZENS: Citizen[] = [
  { id: '1', name: 'Abubakar Sani', nnu: 'SN-88472', status: 'active', agency: 'NIA HQ' },
  { id: '2', name: 'Chioma Okafor', nnu: 'SN-99103', status: 'active', agency: 'Lagos State' },
  { id: '3', name: 'Emeka Nwosu', nnu: 'SN-77231', status: 'suspended', agency: 'FCT Admin' },
  { id: '4', name: 'Fatima Usman', nnu: 'SN-44567', status: 'pending', agency: 'Kano State' },
  { id: '5', name: 'Gabriel Okonkwo', nnu: 'SN-33891', status: 'active', agency: 'NIA HQ' },
  { id: '6', name: 'Hauwa Ibrahim', nnu: 'SN-22904', status: 'revoked', agency: 'Rivers State' },
  { id: '7', name: 'Ifeanyi Eze', nnu: 'SN-11782', status: 'active', agency: 'Enugu State' },
  { id: '8', name: 'Joy Akpan', nnu: 'SN-00673', status: 'pending', agency: 'Akwa Ibom' },
  { id: '9', name: 'Kelechi Obi', nnu: 'SN-55489', status: 'active', agency: 'NIA HQ' },
  { id: '10', name: 'Lami John', nnu: 'SN-44321', status: 'active', agency: 'Plateau State' },
  { id: '11', name: 'Musa Bello', nnu: 'SN-33210', status: 'suspended', agency: 'Kaduna State' },
  { id: '12', name: 'Ngozi Eze', nnu: 'SN-22109', status: 'active', agency: 'Anambra State' },
];

const AGENCIES = ['All', 'NIA HQ', 'Lagos State', 'FCT Admin', 'Kano State', 'Rivers State', 'Enugu State', 'Akwa Ibom', 'Plateau State', 'Kaduna State', 'Anambra State'];
const STATUSES = ['All', 'active', 'pending', 'suspended', 'revoked'];

export default function CitizensScreen() {
  const theme = useTheme();
  const router = useRouter();
  const [search, setSearch] = useState('');
  const [refreshing, setRefreshing] = useState(false);
  const [showFilter, setShowFilter] = useState(false);
  const [statusFilter, setStatusFilter] = useState<string>('All');
  const [agencyFilter, setAgencyFilter] = useState<string>('All');

  const filtered = MOCK_CITIZENS.filter(c => {
    const matchesSearch = search === '' ||
      c.name.toLowerCase().includes(search.toLowerCase()) ||
      c.nnu.toLowerCase().includes(search.toLowerCase());
    const matchesStatus = statusFilter === 'All' || c.status === statusFilter;
    const matchesAgency = agencyFilter === 'All' || c.agency === agencyFilter;
    return matchesSearch && matchesStatus && matchesAgency;
  });

  const onRefresh = useCallback(() => {
    setRefreshing(true);
    setTimeout(() => setRefreshing(false), 1500);
  }, []);

  return (
    <SafeAreaView style={[styles.safeArea, { backgroundColor: theme.background }]}>
      <FlatList
        data={filtered}
        keyExtractor={(item) => item.id}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
        contentContainerStyle={styles.listContent}
        ListHeaderComponent={
          <ThemedView>
            <View style={styles.header}>
              <ThemedText type="subtitle" style={{ fontSize: 26, lineHeight: 32 }}>Citizens</ThemedText>
              <ThemedText themeColor="textSecondary" style={{ fontSize: 13 }}>{filtered.length} records</ThemedText>
            </View>
            <View style={styles.searchWrapper}>
              <SearchBar
                placeholder="Search by name or NNU..."
                value={search}
                onChangeText={setSearch}
                onFilterPress={() => setShowFilter(true)}
              />
            </View>
          </ThemedView>
        }
        renderItem={({ item }) => (
          <View style={styles.cardWrapper}>
            <CitizenCard
              name={item.name}
              nnu={item.nnu}
              status={item.status}
              agency={item.agency}
              onPress={() => router.push(`/citizen/${item.id}`)}
            />
          </View>
        )}
        ListEmptyComponent={
          <ThemedView style={styles.empty}>
            <ThemedText themeColor="textSecondary" style={{ textAlign: 'center' }}>
              {search || statusFilter !== 'All' || agencyFilter !== 'All'
                ? 'No citizens match your filters.'
                : 'No citizens found.'}
            </ThemedText>
          </ThemedView>
        }
      />

      <Modal visible={showFilter} animationType="slide" transparent>
        <Pressable style={styles.modalOverlay} onPress={() => setShowFilter(false)}>
          <Pressable onPress={(e) => e.stopPropagation()}>
            <ThemedView type="backgroundElement" style={styles.filterModal}>
              <ThemedText style={{ fontWeight: '600', fontSize: 17, marginBottom: Spacing.three }}>Filter Citizens</ThemedText>

              <ThemedText themeColor="textSecondary" style={{ fontSize: 13, marginBottom: Spacing.two }}>Status</ThemedText>
              <View style={styles.chipRow}>
                {STATUSES.map(s => (
                  <Pressable
                    key={s}
                    onPress={() => setStatusFilter(s)}
                    style={[styles.chip, { backgroundColor: statusFilter === s ? theme.text : theme.backgroundSelected }]}
                  >
                    <ThemedText style={{ color: statusFilter === s ? theme.background : theme.text, fontSize: 13, textTransform: 'capitalize' }}>
                      {s}
                    </ThemedText>
                  </Pressable>
                ))}
              </View>

              <ThemedText themeColor="textSecondary" style={{ fontSize: 13, marginBottom: Spacing.two, marginTop: Spacing.three }}>Agency</ThemedText>
              <View style={styles.chipRow}>
                {AGENCIES.map(a => (
                  <Pressable
                    key={a}
                    onPress={() => setAgencyFilter(a)}
                    style={[styles.chip, { backgroundColor: agencyFilter === a ? theme.text : theme.backgroundSelected }]}
                  >
                    <ThemedText style={{ color: agencyFilter === a ? theme.background : theme.text, fontSize: 13 }}>
                      {a}
                    </ThemedText>
                  </Pressable>
                ))}
              </View>

              <Pressable onPress={() => { setStatusFilter('All'); setAgencyFilter('All'); }} style={styles.resetBtn}>
                <ThemedText style={{ color: '#1565c0', fontWeight: '600', fontSize: 13 }}>Reset Filters</ThemedText>
              </Pressable>
            </ThemedView>
          </Pressable>
        </Pressable>
      </Modal>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: {
    flex: 1,
  },
  listContent: {
    paddingBottom: BottomTabInset + Spacing.three,
    maxWidth: MaxContentWidth,
    width: '100%',
    alignSelf: 'center',
  },
  header: {
    paddingHorizontal: Spacing.four,
    paddingTop: Spacing.three,
    gap: Spacing.half,
  },
  searchWrapper: {
    paddingHorizontal: Spacing.four,
    paddingVertical: Spacing.three,
  },
  cardWrapper: {
    paddingHorizontal: Spacing.four,
    marginBottom: Spacing.two,
  },
  empty: {
    padding: Spacing.six,
    alignItems: 'center',
  },
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0,0,0,0.4)',
    justifyContent: 'flex-end',
  },
  filterModal: {
    borderTopLeftRadius: Spacing.four,
    borderTopRightRadius: Spacing.four,
    padding: Spacing.four,
    maxHeight: '70%',
  },
  chipRow: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: Spacing.two,
  },
  chip: {
    paddingHorizontal: Spacing.three,
    paddingVertical: Spacing.two,
    borderRadius: Spacing.five,
  },
  resetBtn: {
    alignSelf: 'center',
    marginTop: Spacing.four,
    paddingVertical: Spacing.two,
  },
});
