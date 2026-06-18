package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/biar/internal/domain"
)

type WeaponRepository interface {
	Create(ctx context.Context, w *domain.IllicitWeapon) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.IllicitWeapon, error)
	List(ctx context.Context) ([]*domain.IllicitWeapon, error)
	CheckSerial(ctx context.Context, serial string) ([]*domain.IllicitWeapon, error)
	UpsertFromIARMS(ctx context.Context, w *domain.IllicitWeapon) error
	ByGang(ctx context.Context) ([]*domain.WeaponsByGang, error)
	ByOrigin(ctx context.Context) ([]*domain.WeaponsByOrigin, error)
	Routes(ctx context.Context) ([]*domain.TraffickingRoute, error)
}

type InMemoryWeaponRepo struct {
	mu      sync.RWMutex
	weapons map[uuid.UUID]*domain.IllicitWeapon
}

func NewInMemoryWeaponRepo() *InMemoryWeaponRepo {
	return &InMemoryWeaponRepo{
		weapons: make(map[uuid.UUID]*domain.IllicitWeapon),
	}
}

func (r *InMemoryWeaponRepo) Create(ctx context.Context, w *domain.IllicitWeapon) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.weapons[w.WeaponID] = w
	return nil
}

func (r *InMemoryWeaponRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.IllicitWeapon, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, ok := r.weapons[id]
	if !ok {
		return nil, ErrNotFound
	}
	return w, nil
}

func (r *InMemoryWeaponRepo) List(ctx context.Context) ([]*domain.IllicitWeapon, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.IllicitWeapon, 0, len(r.weapons))
	for _, w := range r.weapons {
		result = append(result, w)
	}
	return result, nil
}

func (r *InMemoryWeaponRepo) CheckSerial(ctx context.Context, serial string) ([]*domain.IllicitWeapon, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.IllicitWeapon
	for _, w := range r.weapons {
		if w.SerialNumber != nil && *w.SerialNumber == serial {
			result = append(result, w)
		}
	}
	return result, nil
}

func (r *InMemoryWeaponRepo) UpsertFromIARMS(ctx context.Context, w *domain.IllicitWeapon) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.weapons[w.WeaponID] = w
	return nil
}

func (r *InMemoryWeaponRepo) ByGang(ctx context.Context) ([]*domain.WeaponsByGang, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	gangCount := make(map[string]*domain.WeaponsByGang)
	for _, w := range r.weapons {
		if w.GangID == nil {
			continue
		}
		gid := w.GangID.String()
		entry, ok := gangCount[gid]
		if !ok {
			entry = &domain.WeaponsByGang{
				GangID:      gid,
				GangName:    fmt.Sprintf("Gang %s", gid[:8]),
				WeaponCount: 0,
			}
			gangCount[gid] = entry
		}
		entry.WeaponCount++
	}
	result := make([]*domain.WeaponsByGang, 0, len(gangCount))
	for _, v := range gangCount {
		result = append(result, v)
	}
	return result, nil
}

func (r *InMemoryWeaponRepo) ByOrigin(ctx context.Context) ([]*domain.WeaponsByOrigin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	originCount := make(map[string]int)
	for _, w := range r.weapons {
		if w.OriginCountry == nil {
			continue
		}
		originCount[*w.OriginCountry]++
	}
	result := make([]*domain.WeaponsByOrigin, 0, len(originCount))
	for country, count := range originCount {
		result = append(result, &domain.WeaponsByOrigin{
			OriginCountry: country,
			WeaponCount:   count,
		})
	}
	return result, nil
}

func (r *InMemoryWeaponRepo) Routes(ctx context.Context) ([]*domain.TraffickingRoute, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	routeMap := make(map[string]*domain.TraffickingRoute)
	for _, w := range r.weapons {
		if w.OriginCountry == nil {
			continue
		}
		key := *w.OriginCountry
		entry, ok := routeMap[key]
		if !ok {
			entry = &domain.TraffickingRoute{
				OriginCountry: *w.OriginCountry,
				WeaponCount:   0,
			}
			routeMap[key] = entry
		}
		if w.TransitCountries != nil {
			entry.TransitCountries = w.TransitCountries
		}
		if w.ImportMethod != nil {
			entry.ImportMethod = *w.ImportMethod
		}
		entry.WeaponCount++
	}
	result := make([]*domain.TraffickingRoute, 0, len(routeMap))
	for _, v := range routeMap {
		result = append(result, v)
	}
	return result, nil
}
