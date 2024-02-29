export enum Type {
  Sending = 'sending',
  Receiving = 'receiving'
}

export type Accounting = {
  id: string
  account_number: string
  account_name: string
  iban: string
  address: string
  amount: number
  type: Type
}
