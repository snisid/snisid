import { useState, useCallback, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TextInput,
  TouchableOpacity,
  ActivityIndicator,
  Alert,
  KeyboardAvoidingView,
  Platform,
} from 'react-native';
import { useRouter, Stack } from 'expo-router';

import { useAuth } from '@/lib/auth';
import type { Identity } from '@/lib/api';

interface FormErrors {
  fullName?: string;
  dateOfBirth?: string;
  nationality?: string;
  gender?: string;
}

export default function IdentityEditScreen() {
  const router = useRouter();
  const { api } = useAuth();

  const [identity, setIdentity] = useState<Identity | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');

  const [fullName, setFullName] = useState('');
  const [dateOfBirth, setDateOfBirth] = useState('');
  const [nationality, setNationality] = useState('');
  const [gender, setGender] = useState('');
  const [errors, setErrors] = useState<FormErrors>({});

  const fetchIdentity = useCallback(async () => {
    try {
      setError('');
      const response = await api.getIdentity();
      setIdentity(response.identity);
      setFullName(response.identity.fullName);
      setDateOfBirth(response.identity.dateOfBirth.split('T')[0]);
      setNationality(response.identity.nationality);
      setGender(response.identity.gender);
    } catch (err: any) {
      setError(
        err?.response?.data?.message ?? err?.message ?? 'Failed to load identity'
      );
    } finally {
      setLoading(false);
    }
  }, [api]);

  useEffect(() => {
    fetchIdentity();
  }, [fetchIdentity]);

  const validate = useCallback((): FormErrors => {
    const errs: FormErrors = {};
    if (!fullName.trim() || fullName.trim().length < 3) {
      errs.fullName = 'Full name must be at least 3 characters';
    }
    if (fullName.trim().split(' ').length < 2) {
      errs.fullName = errs.fullName ?? 'Please enter both first and last name';
    }
    if (dateOfBirth) {
      const dob = new Date(dateOfBirth);
      if (isNaN(dob.getTime())) {
        errs.dateOfBirth = 'Invalid date format';
      } else {
        const age = new Date().getFullYear() - dob.getFullYear();
        if (age < 16) errs.dateOfBirth = 'You must be at least 16 years old';
        if (age > 150) errs.dateOfBirth = 'Invalid date of birth';
      }
    }
    if (!nationality.trim()) {
      errs.nationality = 'Nationality is required';
    }
    if (!gender.trim()) {
      errs.gender = 'Gender is required';
    }
    return errs;
  }, [fullName, dateOfBirth, nationality, gender]);

  const hasChanges = useCallback((): boolean => {
    if (!identity) return false;
    return (
      fullName !== identity.fullName ||
      nationality !== identity.nationality ||
      gender !== identity.gender
    );
  }, [identity, fullName, nationality, gender]);

  const handleSubmit = useCallback(() => {
    const validationErrors = validate();
    setErrors(validationErrors);
    if (Object.keys(validationErrors).length > 0) return;

    if (!hasChanges()) {
      Alert.alert('No Changes', 'No fields have been modified.');
      return;
    }

    Alert.alert(
      'Confirm Update',
      'Are you sure you want to update your identity information? This action will be logged and may require verification.',
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Confirm Update',
          style: 'destructive',
          onPress: async () => {
            setSaving(true);
            try {
              const payload: Partial<{ fullName: string; dateOfBirth: string; nationality: string; gender: string }> = {};
              if (fullName !== identity?.fullName) payload.fullName = fullName.trim();
              if (nationality !== identity?.nationality) payload.nationality = nationality.trim();
              if (gender !== identity?.gender) payload.gender = gender.trim();

              await api.updateIdentity(payload);
              Alert.alert('Update Submitted', 'Your identity update request has been submitted for review.', [
                { text: 'OK', onPress: () => router.back() },
              ]);
            } catch (err: any) {
              Alert.alert(
                'Update Failed',
                err?.response?.data?.message ?? err?.message ?? 'Failed to update identity'
              );
            } finally {
              setSaving(false);
            }
          },
        },
      ]
    );
  }, [validate, hasChanges, api, identity, fullName, nationality, gender, router]);

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#0033a0" />
        <Text style={styles.loadingText}>Loading identity data...</Text>
      </View>
    );
  }

  if (error) {
    return (
      <View style={styles.center}>
        <Text style={styles.errorIcon}>⚠️</Text>
        <Text style={styles.errorText}>{error}</Text>
        <TouchableOpacity style={styles.retryButton} onPress={fetchIdentity}>
          <Text style={styles.retryText}>Retry</Text>
        </TouchableOpacity>
      </View>
    );
  }

  return (
    <KeyboardAvoidingView
      style={styles.flex}
      behavior={Platform.OS === 'ios' ? 'padding' : undefined}>
      <Stack.Screen
        options={{
          headerShown: true,
          headerTitle: 'Update Identity',
          headerStyle: { backgroundColor: '#F5F5F7' },
          headerTintColor: '#0033a0',
        }}
      />
      <ScrollView style={styles.container} contentContainerStyle={styles.content}>
        <Text style={styles.pageTitle}>Edit Identity</Text>
        <Text style={styles.pageSubtitle}>
          Update your personal information. Changes require verification.
        </Text>

        <View style={styles.formCard}>
          <View style={styles.field}>
            <Text style={styles.label}>Full Name</Text>
            <TextInput
              style={[styles.input, errors.fullName && styles.inputError]}
              value={fullName}
              onChangeText={(t) => {
                setFullName(t);
                if (errors.fullName) setErrors((prev) => ({ ...prev, fullName: undefined }));
              }}
              placeholder="e.g. Thabo Mbeki"
              placeholderTextColor="#999"
              autoCapitalize="words"
            />
            {errors.fullName && <Text style={styles.fieldError}>{errors.fullName}</Text>}
          </View>

          <View style={styles.field}>
            <Text style={styles.label}>Date of Birth</Text>
            <TextInput
              style={[styles.input, errors.dateOfBirth && styles.inputError]}
              value={dateOfBirth}
              onChangeText={(t) => {
                setDateOfBirth(t);
                if (errors.dateOfBirth) setErrors((prev) => ({ ...prev, dateOfBirth: undefined }));
              }}
              placeholder="YYYY-MM-DD"
              placeholderTextColor="#999"
            />
            {errors.dateOfBirth && (
              <Text style={styles.fieldError}>{errors.dateOfBirth}</Text>
            )}
          </View>

          <View style={styles.field}>
            <Text style={styles.label}>Nationality</Text>
            <TextInput
              style={[styles.input, errors.nationality && styles.inputError]}
              value={nationality}
              onChangeText={(t) => {
                setNationality(t);
                if (errors.nationality) setErrors((prev) => ({ ...prev, nationality: undefined }));
              }}
              placeholder="e.g. South African"
              placeholderTextColor="#999"
              autoCapitalize="words"
            />
            {errors.nationality && (
              <Text style={styles.fieldError}>{errors.nationality}</Text>
            )}
          </View>

          <View style={styles.field}>
            <Text style={styles.label}>Gender</Text>
            <View style={styles.genderRow}>
              {['Male', 'Female', 'Other'].map((opt) => (
                <TouchableOpacity
                  key={opt}
                  style={[
                    styles.genderOption,
                    gender === opt && styles.genderOptionActive,
                  ]}
                  onPress={() => {
                    setGender(opt);
                    if (errors.gender) setErrors((prev) => ({ ...prev, gender: undefined }));
                  }}>
                  <Text
                    style={[
                      styles.genderText,
                      gender === opt && styles.genderTextActive,
                    ]}>
                    {opt}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>
            {errors.gender && <Text style={styles.fieldError}>{errors.gender}</Text>}
          </View>

          {identity && (
            <View style={styles.currentValues}>
              <Text style={styles.currentLabel}>Current Values</Text>
              <Text style={styles.currentText}>NNU: {identity.nnu}</Text>
              <Text style={styles.currentText}>Name: {identity.fullName}</Text>
              <Text style={styles.currentText}>
                DOB: {new Date(identity.dateOfBirth).toLocaleDateString()}
              </Text>
              <Text style={styles.currentText}>
                Status: {identity.status.toUpperCase()}
              </Text>
            </View>
          )}
        </View>

        <TouchableOpacity
          style={[styles.submitButton, saving && styles.submitButtonDisabled]}
          onPress={handleSubmit}
          disabled={saving}>
          {saving ? (
            <ActivityIndicator size="small" color="#fff" />
          ) : (
            <Text style={styles.submitText}>Submit Update Request</Text>
          )}
        </TouchableOpacity>
      </ScrollView>
    </KeyboardAvoidingView>
  );
}

const styles = StyleSheet.create({
  flex: {
    flex: 1,
  },
  container: {
    flex: 1,
    backgroundColor: '#F5F5F7',
  },
  content: {
    paddingTop: 24,
    paddingBottom: 40,
    gap: 20,
  },
  center: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 24,
    backgroundColor: '#F5F5F7',
  },
  pageTitle: {
    fontSize: 28,
    fontWeight: '700',
    color: '#000',
    paddingHorizontal: 20,
  },
  pageSubtitle: {
    fontSize: 14,
    color: '#666',
    paddingHorizontal: 20,
    lineHeight: 20,
  },
  loadingText: {
    marginTop: 12,
    fontSize: 14,
    color: '#666',
  },
  errorIcon: {
    fontSize: 48,
    marginBottom: 12,
  },
  errorText: {
    fontSize: 14,
    color: '#991B1B',
    textAlign: 'center',
    marginBottom: 16,
  },
  retryButton: {
    backgroundColor: '#0033a0',
    paddingHorizontal: 24,
    paddingVertical: 12,
    borderRadius: 8,
  },
  retryText: {
    color: '#fff',
    fontWeight: '600',
  },
  formCard: {
    backgroundColor: '#fff',
    borderRadius: 16,
    marginHorizontal: 16,
    padding: 20,
    gap: 20,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.05,
    shadowRadius: 8,
    elevation: 2,
  },
  field: {
    gap: 6,
  },
  label: {
    fontSize: 13,
    fontWeight: '600',
    color: '#555',
    textTransform: 'uppercase',
    letterSpacing: 0.3,
  },
  input: {
    borderWidth: 1,
    borderColor: '#E5E5EA',
    borderRadius: 10,
    padding: 14,
    fontSize: 16,
    color: '#000',
    backgroundColor: '#FAFAFA',
  },
  inputError: {
    borderColor: '#DC2626',
    backgroundColor: '#FEF2F2',
  },
  fieldError: {
    fontSize: 12,
    color: '#DC2626',
  },
  genderRow: {
    flexDirection: 'row',
    gap: 10,
  },
  genderOption: {
    flex: 1,
    paddingVertical: 12,
    borderRadius: 10,
    borderWidth: 1,
    borderColor: '#E5E5EA',
    alignItems: 'center',
    backgroundColor: '#FAFAFA',
  },
  genderOptionActive: {
    backgroundColor: '#0033a0',
    borderColor: '#0033a0',
  },
  genderText: {
    fontSize: 15,
    fontWeight: '600',
    color: '#555',
  },
  genderTextActive: {
    color: '#fff',
  },
  currentValues: {
    backgroundColor: '#F0F4FF',
    borderRadius: 10,
    padding: 14,
    gap: 4,
  },
  currentLabel: {
    fontSize: 12,
    fontWeight: '700',
    color: '#0033a0',
    textTransform: 'uppercase',
    letterSpacing: 0.3,
    marginBottom: 4,
  },
  currentText: {
    fontSize: 13,
    color: '#555',
    fontFamily: Platform.OS === 'ios' ? 'monospace' : 'monospace',
  },
  submitButton: {
    backgroundColor: '#0033a0',
    marginHorizontal: 16,
    borderRadius: 14,
    paddingVertical: 16,
    alignItems: 'center',
  },
  submitButtonDisabled: {
    opacity: 0.6,
  },
  submitText: {
    color: '#fff',
    fontSize: 17,
    fontWeight: '600',
  },
});
