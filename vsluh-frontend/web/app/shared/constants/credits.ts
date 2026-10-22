/** Пакеты кредитов. Один кредит — один разбор речи. */
export interface CreditPack {
  credits: number
  price: number
}

export const CREDIT_PACKS: CreditPack[] = [
  { credits: 10, price: 199 },
  { credits: 50, price: 799 },
  { credits: 150, price: 1990 },
]
