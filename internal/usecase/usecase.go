package usecase

type UseCase struct {
	Auth Auth
}

func New(auth Auth) *UseCase {
	return &UseCase{
		Auth: auth,
	}
}
