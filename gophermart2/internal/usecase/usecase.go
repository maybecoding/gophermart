package usecase

type UseCase struct {
	Auth  Auth
	Order Order
	Bonus Bonus
}

func New(auth Auth, order Order, bonus Bonus) *UseCase {
	return &UseCase{
		Auth:  auth,
		Order: order,
		Bonus: bonus,
	}
}
