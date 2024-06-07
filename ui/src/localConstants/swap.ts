
// All events that can be emitted within a swap context. It can be for both the
// initiator and the participant.
export enum SwapEvents {
  SInit = 'SInit',
  SInitDone = 'SInitDone',
  SLockedEGLD = 'SLockedEGLD',
  SLockeedBTC = 'SLockedBTC',
  SRefund = 'SRefund',
  SClaimed = 'SClaimed',
  SOk = 'SOk',
  SFailed = 'SFailed',
};
