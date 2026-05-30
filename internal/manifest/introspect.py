import argparse
import inspect
import importlib
import json
import sys


DOCUMENTATION_BASE = "https://www.tensorflow.org/api_docs/python"


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--module", default="tensorflow")
    parser.add_argument("--root", default="tf")
    parser.add_argument("--max-depth", type=int, default=3)
    args = parser.parse_args()

    module = importlib.import_module(args.module)
    manifest = {
        "schema_version": 1,
        "source_module": args.module,
        "source_version": getattr(module, "__version__", ""),
        "generated_by": "tftf-manifest",
        "root": args.root,
        "documentation_base": DOCUMENTATION_BASE,
        "symbols": collect_symbols(module, args.root, args.max_depth),
    }

    json.dump(manifest, sys.stdout, indent=2, sort_keys=False)
    sys.stdout.write("\n")


def collect_symbols(root_obj, root_name, max_depth):
    symbols = {}
    queue = [(root_name, root_obj, 0)]
    seen_objects = set()

    while queue:
        path, obj, depth = queue.pop(0)
        kind = classify(obj)
        children = []

        if should_descend(obj, depth, max_depth):
            object_id = id(obj)
            if object_id not in seen_objects:
                seen_objects.add(object_id)
                for name in sorted(public_names(obj)):
                    child = safe_getattr(obj, name)
                    if child is None:
                        continue
                    child_path = f"{path}.{name}"
                    children.append(child_path)
                    queue.append((child_path, child, depth + 1))

        symbols[path] = {
            "path": path,
            "kind": kind,
            "module": getattr(obj, "__module__", ""),
            "signature": signature_for(obj),
            "doc_url": doc_url(path),
            "children": children,
        }

    return [clean(symbols[path]) for path in sorted(symbols)]


def public_names(obj):
    try:
        names = dir(obj)
    except Exception:
        return []
    return [name for name in names if not name.startswith("_")]


def safe_getattr(obj, name):
    try:
        return getattr(obj, name)
    except Exception:
        return None


def should_descend(obj, depth, max_depth):
    if depth >= max_depth:
        return False
    return inspect.ismodule(obj) or inspect.isclass(obj)


def classify(obj):
    if inspect.ismodule(obj):
        return "module"
    if inspect.isclass(obj):
        return "class"
    if inspect.isfunction(obj) or inspect.ismethod(obj) or inspect.isbuiltin(obj):
        return "function"
    if callable(obj):
        return "callable"
    return "value"


def signature_for(obj):
    if not callable(obj):
        return ""
    try:
        return str(inspect.signature(obj))
    except Exception:
        return ""


def doc_url(path):
    return f"{DOCUMENTATION_BASE}/{path.replace('.', '/')}"


def clean(symbol):
    return {key: value for key, value in symbol.items() if value not in ("", [], None)}


if __name__ == "__main__":
    main()
