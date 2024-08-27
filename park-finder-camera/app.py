from src import create_app


def main():
    app = create_app()
    app.run(host="0.0.0.0", port=3000, debug=True)


if __name__ == "__main__":
    main()
    