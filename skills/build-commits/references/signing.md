# Commit signing

Use only keys the user already owns and has configured. Never generate, import or upload keys without asking.

## Check current state
```bash
git config --get commit.gpgsign      # true = sign every commit
git config --get gpg.format          # openpgp (default) | ssh | x509
git config --get user.signingkey     # key id, or path to SSH public key
git log -1 --show-signature          # verify the last commit
git log --format='%h %G? %GS %s' origin/<base>..HEAD   # G=good, U=good untrusted, N=none, B=bad, E=cannot check
```

## SSH signing (simplest, git 2.34+)
```bash
git config gpg.format ssh
git config user.signingkey ~/.ssh/id_ed25519.pub
git config commit.gpgsign true
```
Local verification needs an allowed-signers file: `git config gpg.ssh.allowedSignersFile ~/.ssh/allowed_signers` with lines `<email> <public key>`. Hosting platforms verify against keys the user registered as signing keys (not just authentication keys).

## GPG signing
```bash
gpg --list-secret-keys --keyid-format=long
git config user.signingkey <LONGKEYID>
git config commit.gpgsign true
```
The key's UID email must match `user.email` and be verified on the hosting platform, or the commit shows "Unverified".

## Common failures
| Symptom | Cause | Fix |
|---|---|---|
| `error: gpg failed to sign the data` | no agent / pinentry in a non-interactive shell | set `GPG_TTY=$(tty)`, or ask the user to sign; do not disable signing |
| Commit shows Unverified | email mismatch, key not uploaded, expired key | align `user.email` with the key, ask the user to register the key |
| Rebased commits lost signatures | rebase re-creates commits | keep `commit.gpgsign=true` during rebase, or `git rebase --exec 'git commit --amend --no-edit -S' <base>` |
| Squash merge on platform | platform signs the squash commit with its own key | expected; branch commits still need signing if a ruleset requires it |
| `--no-gpg-sign` suggested | bypass | never use it to get past a requirement |

## Tags
Signed tags: `git tag -s vX.Y.Z -m "vX.Y.Z"`; verify with `git tag -v vX.Y.Z`. See `ship-release`.

## Policy
- If a ruleset requires signed commits, unsigned commits cannot be merged; fix signing before pushing rather than rewriting later.
- Signing proves who committed, not that the code is safe; it complements review, not replaces it.
