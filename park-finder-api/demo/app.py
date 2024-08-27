from src import create_app


def main():
    app = create_app()
    app.run(host="0.0.0.0", port=4500, debug=True)


if __name__ == "__main__":
    main()
    