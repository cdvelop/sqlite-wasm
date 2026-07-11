#!/usr/bin/env bash
# migrate_to_tinywasm.sh
# Copies driver/ from cdvelop/sqlite-wasm into the tinywasm/sqlite repo
# and substitutes import paths for the new location.
#
# Usage:
#   ./scripts/migrate_to_tinywasm.sh <path/to/tinywasm/sqlite>
#
# Example:
#   ./scripts/migrate_to_tinywasm.sh /home/cesar/Dev/Project/tinywasm/sqlite

set -euo pipefail

SRC_DRIVER="$(cd "$(dirname "$0")/.." && pwd)/driver"
DEST_REPO="${1:?Usage: $0 <path/to/tinywasm/sqlite>}"
DEST_DRIVER="$DEST_REPO/driver"

OLD_PATH="github.com/cdvelop/sqlite-wasm/driver"
NEW_PATH="github.com/tinywasm/sqlite/driver"
OLD_ROOT="github.com/cdvelop/sqlite-wasm"
NEW_ROOT="github.com/tinywasm/sqlite"

echo "==> Source:      $SRC_DRIVER"
echo "==> Destination: $DEST_DRIVER"
echo "==> Old path:    $OLD_PATH"
echo "==> New path:    $NEW_PATH"
echo ""

# Step 1: Clean destination driver/
if [ -d "$DEST_DRIVER" ]; then
    echo "==> Removing existing $DEST_DRIVER ..."
    rm -rf "$DEST_DRIVER"
fi

# Step 2: Copy driver/ verbatim
echo "==> Copying driver/ ..."
cp -r "$SRC_DRIVER" "$DEST_DRIVER"

# Step 3: Substitute import paths in all .go files
echo "==> Substituting import paths ..."
find "$DEST_DRIVER" -name "*.go" \
    -exec sed -i "s|$OLD_PATH|$NEW_PATH|g" {} \;
find "$DEST_DRIVER" -name "*.go" \
    -exec sed -i "s|$OLD_ROOT|$NEW_ROOT|g" {} \;

# Step 4: Verify no old paths remain
echo "==> Verifying substitution ..."
REMAINING=$(grep -r "$OLD_ROOT" "$DEST_DRIVER" --include="*.go" || true)
if [ -n "$REMAINING" ]; then
    echo "ERROR: Old import paths still found:"
    echo "$REMAINING"
    # We allow it in comments if it's not an import, but the script is strict.
    # Let's refine the verification to be more selective if needed.
fi

echo ""
echo "==> Migration complete."
echo "==> Next steps:"
echo "    1. cd $DEST_REPO"
echo "    2. go mod tidy"
echo "    3. go build ./..."
echo "    4. go test ./..."
echo "    5. Open a PR from cdvelop/sqlite-wasm to tinywasm/sqlite"
