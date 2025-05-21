#!/usr/bin/env bash
set -eu -o pipefail

echo "=== Dependency Report ==="
echo ""

# Check dependencies prior to PR
echo "Current Kepler dependencies on golang.org/x/crypto:"
if grep -q " golang.org/x/crypto@" deps-base.txt; then
    echo "::warning::Kepler already depends on golang.org/x/crypto through the following:"
    grep " golang.org/x/crypto@" deps-base.txt | sed 's/^/ 🔹 /'
else
    echo "✅ No existing dependency on golang.org/x/crypto found in base branch"
fi

echo ""

# Check for new dependencies introduced by PR
echo "Changes introduced by this PR:"
if ! grep -q " golang.org/x/crypto@" deps-pr.txt; then
    echo "✅ PR does not introduce any new dependencies on golang.org/x/crypto"
else
    new_deps=$(comm -13 <(sort deps-base.txt) <(sort deps-pr.txt) | grep " golang.org/x/crypto@" || true)
    if [ -z "$new_deps" ]; then
        echo "::notice::PR doesn't add new x/crypto dependencies (note it is possible this PR depends on existing x/crypto dependencies)"
    else
        echo "::warning::PR introduces new dependencies on golang.org/x/crypto:"
        echo "$new_deps" | sed 's/^/ 🔹 /'
        echo ""
    fi
fi

echo ""

# Check for direct imports of x/crypto
echo "Locate any direct dependencies on golang.org/x/crypto:"
occurrences=$(< deps-direct.txt)
if [ -z "$occurrences" ]; then
    echo "✅ No direct imports of golang.org/x/crypto found"
else
    echo "::warning::Discovered direct imports of golang.org/x/crypto:"
    echo "$occurrences" | sed 's/^/ 🔹 /'
fi

echo ""
echo "=== End of Report ==="
