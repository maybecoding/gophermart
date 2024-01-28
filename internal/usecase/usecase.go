package usecase

type UseCase struct {
	Auth  Auth
	Order Order
}

func New(auth Auth, order Order) *UseCase {
	return &UseCase{
		Auth:  auth,
		Order: order,
	}
}
