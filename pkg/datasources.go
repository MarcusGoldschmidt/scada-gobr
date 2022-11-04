package pkg

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/datasources"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
)

func (s *Scadagobr) LoadDataSourcesRuntimeManager(ctx context.Context) ([]datasources.DataSourceRuntimeManager, error) {
	ds, err := s.dataSourcePersistence.GetDataSources(ctx)

	if err != nil {
		return nil, err
	}

	dsrm := make([]datasources.DataSourceRuntimeManager, len(ds))

	for i, d := range ds {
		manager, err := DataSourceToRuntimeManager(s, d)
		if err != nil {
			return nil, err
		}

		dsrm[i] = manager
	}

	return dsrm, nil
}

func (s *Scadagobr) UpdateDataSource(ctx context.Context, id shared.CommonId) error {
	ds, err := s.dataSourcePersistence.GetDataSourceById(ctx, id)

	if err != nil {
		return err
	}

	manager, err := DataSourceToRuntimeManager(s, ds)

	err = s.RuntimeManager.UpdateDataSource(ctx, manager)
	if err != nil {
		return err
	}

	return nil
}
